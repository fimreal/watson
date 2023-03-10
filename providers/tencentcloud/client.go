package tencentcloud

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/fimreal/goutils/ezap"
	"github.com/fimreal/watson/providers"
	tc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	th "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
	tp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func (p *Provider) doRequest(action string, payload any) ([]byte, error) {

	ezap.Debugf("Request tencentcloud API, action: %s, payload: %+v", action, payload)

	cpf := tp.NewClientProfile()
	cpf.HttpProfile.RootDomain = "tencentcloudapi.com"
	credential := tc.NewCredential(p.SecretId, p.SecretKey)
	client := tc.NewCommonClient(credential, "", cpf)
	request := th.NewCommonRequest("dnspod", "2021-03-23", action)
	request.SetActionParameters(payload)

	response := th.NewCommonResponse()
	if err := client.Send(request, response); err != nil {
		return nil, err
	}

	return response.GetBody(), nil
}

// 获取某个域名下的解析记录列表
// 默认接口请求频率限制：100次/秒。
// ref. https://cloud.tencent.com/document/api/1427/80523
func (p *Provider) GetRecordList(zone string) (list []providers.Record, err error) {

	payload := map[string]any{
		"Domain": strings.Trim(zone, "."),
	}

	resp, err := p.doRequest("DescribeRecordList", payload)
	if err != nil {
		return
	}

	data := DescribeRecordListResponse{}
	if err = json.Unmarshal(resp, &data); err != nil {
		return
	}

	for _, record := range data.Response.RecordList {
		list = append(list, providers.Record{
			ID:    strconv.Itoa(record.RecordId),
			Type:  record.Type,
			Name:  record.Name,
			Value: record.Value,
			TTL:   time.Duration(record.TTL),
		})
	}
	return
}

// dnspod 比较蠢的地方在于，获取单个记录需要首先调用 DescribeDomainList 拿到所有记录，获取到要查询的 RecordId，回来再重复查询。
// 这里还是用 DescribeDomainList 取出来，直接筛选。
// ref. https://cloud.tencent.com/document/api/1427/56168
func (p *Provider) GetRecord(zone string, record providers.Record) (ret providers.Record, err error) {
	list, err := p.GetRecordList(zone)
	if err != nil {
		return
	}
	for _, r := range list {
		if r.Name == record.Name {
			return r, nil
		}
	}
	return
}

func (p *Provider) createRecord(zone string, record providers.Record) error {

	payload := map[string]any{
		"Domain":     strings.Trim(zone, "."),
		"SubDomain":  record.Name,
		"RecordType": record.Type,
		"RecordLine": "默认",
		"Value":      record.Value,
	}

	_, err := p.doRequest("CreateRecord", payload)
	return err

}

func (p *Provider) modifyRecord(zone string, record providers.Record) error {

	recordId, _ := strconv.Atoi(record.ID)
	payload := map[string]any{
		"Domain":     strings.Trim(zone, "."),
		"SubDomain":  record.Name,
		"RecordType": record.Type,
		"RecordLine": "默认",
		"Value":      record.Value,
		"RecordId":   recordId,
	}
	_, err := p.doRequest("ModifyRecord", payload)
	return err
}

// 若未取到记录则创建，否则判断是否需要修改然后执行
func (p *Provider) SetRecord(zone string, record providers.Record) error {
	r, err := p.GetRecord(zone, record)
	if err != nil {
		return err
	}
	record.ID = r.ID

	if r.Name == "" {
		return p.createRecord(zone, record)
	}
	if r == record {
		return nil
	}

	ezap.Debugf("%s: %s => %s: %s", r.Name, r.Value, record.Name, record.Value)
	return p.modifyRecord(zone, record)
}

func (p *Provider) DeleteRecord(zone string, record providers.Record) error {
	r, err := p.GetRecord(zone, record)
	if err != nil {
		return err
	}
	if r.Name == "" {
		return nil
	}
	recordId, _ := strconv.Atoi(r.ID)
	payload := map[string]any{
		"Domain":   strings.Trim(zone, "."),
		"RecordId": recordId,
	}
	_, err = p.doRequest("DeleteRecord", payload)
	return err
}
