package chinadep

import (
	"testing"
	"time"

	. "github.com/yuhu-tech/chinadep-go-sdk/internal/util/ptrutil"
)

var (
	appId     = ""
	appSecret = ""
)

func skipTest() bool {
	return appId == "" || appSecret == ""
}

func TestToken(t *testing.T) {
	if skipTest() {
		t.Skip("skipping test due to empty appId or appSecret")
	}
	client := NewClient(appId, appSecret)
	if err := client.ApplyToken(); err != nil {
		t.Fatalf("%v", err)
	}
	if err := client.RefreshToken(); err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%#v", client)
}

func TestAssetsRegister(t *testing.T) {
	if skipTest() {
		t.Skip("skipping test due to empty appId or appSecret")
	}
	type args struct {
		req []*AssetsRegisterRequest
	}
	var tests = []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "fake data from PM",
			args: args{
				req: []*AssetsRegisterRequest{
					{
						MetaverseAssetType: MetaverseAssetTypeDigitalAsset,
						ChainId:            "Meta5y3xbhnu7daq",
						ContractAddr:       "0x3461B67661FE2f9Be3576Fb9a0d1E50933708231",
						SeriesId:           GetStrPtr("0x3461B67661FE2f9Be3576Fb9a0d1E50933708231"),
						SeriesName:         GetStrPtr("0328-01 for Chinadep"),
						SeriesBizType:      GetIntPtr(3),
						SeriesCoverImgUrl:  GetStrPtr("https://providertest.jingkuang.info/server/v1/images/file/system/NDWLnK2E9wKN.png"),
						SeriesDesc:         nil,
						SeriesMediaType:    GetIntPtr(1),
						CreateTime:         GetInt64Ptr(time.Now().Unix()),
						MintNumber:         GetInt64Ptr(10),
						MetadataUrl:        nil,
						IssuerIdType:       GetIntPtr(2),
						IssuerId:           GetStrPtr("111111111111111111"),
						IssuerName:         GetStrPtr("YK测试环境公司"),
						IpIdType:           GetIntPtr(2),
						IpId:               GetStrPtr("222222222222222"),
						IpName:             GetStrPtr("Artist-YK"),
						CirculationInfo: AssetsRegisterRequestCirculationInfo{
							PlatformSelling: GetIntPtr(10),
							Receive:         nil,
							Airdrop:         nil,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	client := NewClient(appId, appSecret)
	if err := client.ApplyToken(); err != nil {
		t.Fatalf("%v", err)
	}
	if err := client.RefreshToken(); err != nil {
		t.Fatalf("%v", err)
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := client.AssetsRegister(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("AssetsRegister() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
