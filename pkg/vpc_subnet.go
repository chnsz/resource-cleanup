package rc

import "github.com/chnsz/resource-cleanup/helper/httphelper"

type VpcSubnet struct {
}

func (*VpcSubnet) QueryIds() []string {
	data, err := httphelper.New("vpc").
		Method("GET").URI("/v1/{project_id}/subnets").
		Request().
		Result()
	if err != nil {
		return nil
	}

	ids := make([]string, 0)
	for _, v := range data.Get("subnets").Array() {
		if v.Get("id").String() == "subnet-default" {
			continue
		}

		ids = append(ids, v.Get("id").String())
	}

	return []string{"5ae0a4d4-83a8-48e4-8d69-a13562ce95bc"}
}

func (*VpcSubnet) GetName() string {
	return "huaweicloud_vpc_subnet"
}
