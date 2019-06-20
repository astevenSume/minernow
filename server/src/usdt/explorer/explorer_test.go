package explorer

import "testing"

func TestExplorer_Balances(t *testing.T) {
	addr1, addr2 := "1FoWyxwPXuj4C6abqwhjDWdz6D4PZgYRjA", "1FoWyxwPXuj4C6abqwhjDWdz6D4PZgYRjA"
	resp, err := NewExplorer().Balances([]string{addr1, addr2})
	if err != nil {
		t.Fatalf("%v", err)
	}

	if _, ok := resp[addr1]; !ok {
		t.Fatalf("address %s info no found", addr1)
	}

	if _, ok := resp[addr2]; !ok {
		t.Fatalf("address %s info no found", addr2)
	}

	t.Logf("%v", resp)
}

func TestExplorer_AddrDetails(t *testing.T) {
	addr1 := "1FoWyxwPXuj4C6abqwhjDWdz6D4PZgYRjA"
	resp, err := NewExplorer().AddrDetails(addr1)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if v, ok := resp["address"]; !ok {
		t.Fatalf("address %s info no found", addr1)
	} else {
		if v.(string) != addr1 {
			t.Fatalf("resp.address %v no equals to %s", v, addr1)
		}
	}
}
