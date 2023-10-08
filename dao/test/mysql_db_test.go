package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/uyouii/oauth/dao/mysql_db"
)

var testDb = mysql_db.GetNewOauthDb(&mysql_db.MysqlDbConfig{
	User:         "root",
	Password:     "asdfgh",
	DatabaseName: "oauth_db",
})

func TestGenPartner(t *testing.T) {
	ctx := context.Background()

	partner, _ := testDb.GenPartner(ctx, "shopee_alert", 24*3600)
	fmt.Printf("partner: %+v\n", partner)
}

func TestGetPartner(t *testing.T) {
	ctx := context.Background()

	partner, _ := testDb.GetPartnerByKey(ctx, "cd66575d0800490c8f75ff9a13d78789")
	fmt.Printf("partner: %+v\n", partner)
}

func TestDeleteParnter(t *testing.T) {
	ctx := context.Background()
	err := testDb.DeletePartner(ctx, "cd66575d0800490c8f75ff9a13d78789")
	fmt.Printf("res: %+v\n", err)
}

func TestGenOauthToken(t *testing.T) {
	ctx := context.Background()
	token, _ := testDb.GenOauthToken(ctx, "089174a3aabb4c8894206aba07f608c2")
	fmt.Printf("token: %+v", token)
}

func TestGetTokenInfo(t *testing.T) {
	ctx := context.Background()
	token, _ := testDb.GetTokenInfo(ctx, "dc896f92004b4257965afaf58dc866ac05d1dfb44742412c8db7730a13fa4819")
	fmt.Printf("token: %+v", token)
}

func TestDeleteExpireTokens(t *testing.T) {
	ctx := context.Background()
	testDb.DeleteExpireTokens(ctx)
}
