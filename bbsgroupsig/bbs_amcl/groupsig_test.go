package bbsgroupsig

import (
	"github.com/milagro-crypto/amcl/version3/go/amcl"
	"testing"
	"time"
)

//测试群签名
func TestGroupSig(t *testing.T)  {

	//测试群签名
	t1 := time.Now()

	rng, err := GetRand()
	if err != nil  {
		t.Fatalf("Error getting rng: \"%s\"", err)
		return
	}

	//生成群公钥
	groupKey, err := NewGroupKey(rng)
	if err != nil {
		t.Fatalf("Group key generation should have succeeded but gave error \"%s\"", err)
		return
	}

	phase1 := time.Since(t1)
	println("Group key generation takes:", phase1/1e6,"ms")

	//生成用户私钥
	t2 := time.Now()

	userKey, err := Registration(groupKey.Gmsk,rng)
	if err != nil {
		t.Fatalf("User key generation should have succeeded but gave error \"%s\"", err)
		return
	}

	phase2 := time.Since(t2)
	println("User key generation takes:", phase2/1e6,"ms")

	//群签名
	t3 := time.Now()

	message := []byte{0,1,2,3,4,5}
	groupSig, err := NewGroupSig(rng,groupKey.GPK,userKey,message)
	if err != nil {
		t.Fatalf("Group signature should be valid but gave error \"%s\"", err)
		return
	}

	phase3 := time.Since(t3)
	println("Creating a group signature takes:", phase3/1e6,"ms")

	//群签名验证
	t4 := time.Now()

	err = GroupVerify(groupSig,groupKey.GPK,message)
	if err != nil {
		t.Fatalf("Group signature should be valid but verification returned error \"%s\"", err)
		return
	}

	phase4 := time.Since(t4)
	println("Verifying a group signature takes:", phase4/1e6,"ms")

	//身份追踪测试
	t5 := time.Now()

	k,err := Tracing(groupSig,groupKey.TK)
	if err != nil {
		t.Fatalf("Tracing should be valid but gave error \"%s\"", err)
		return
	}

	//测试解密出用户私钥是否为原来签名的私钥
	k1 := EcpFromProto(k.Ax)
	uk2 := EcpFromProto(userKey.UK2)
	if !k1.Equals(uk2) {
		t.Fatalf("Tracing shoule be invalid: %s",err)
		return
	}

	phase5 := time.Since(t5)
	println("Tracing takes:", phase5/1e6,"ms")

}

func BenchmarkNewGroupSig(b *testing.B) {
	rng := amcl.NewRAND()
	groupKey, _ := NewGroupKey(rng)
	userKey, _ := Registration(groupKey.Gmsk,rng)
	message := []byte{0,1,2,3,4,5}

	for i:=0;i<b.N;i++ {
		_, _ = NewGroupSig(rng,groupKey.GPK,userKey,message)
	}
}

func BenchmarkGroupVerify(b *testing.B) {
	rng := amcl.NewRAND()
	groupKey, _ := NewGroupKey(rng)
	userKey, _ := Registration(groupKey.Gmsk,rng)
	message := []byte{0,1,2,3,4,5}
	groupsig, _ := NewGroupSig(rng,groupKey.GPK,userKey,message)

	for i:=0;i<b.N;i++ {
		_ = GroupVerify(groupsig,groupKey.GPK,message)
	}
}