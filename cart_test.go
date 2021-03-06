package cart

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "flag"
  "github.com/golang/glog"
)

func TestCheckEnv(t *testing.T) {
  c := New("memory", 1000, "", 100)
  res := c.Mode()
  assert.Equal(t, "memory", res)
  
  c.Close()
}

func TestLoadMem(t *testing.T) {
  c := New("memory", 1000, "", 100)
  res := c.Mode()
  assert.Equal(t, "memory", res)
  
  assert.Nil(t, c.GetCart("1"))
  cart1 := c.AppendSKU("", "product1", "", 1)
  
  cart2 := c.GetCart(cart1.ID)
  assert.NotNil(t, cart2)
  assert.Equal(t, int(1), len(cart2.Items))
  assert.Equal(t, float32(1), cart2.Items[0].Count)
  
  cart1 = c.AppendSKU(cart1.ID, "product1", "", 1)
  cart2 = c.GetCart(cart1.ID)
  assert.Equal(t, float32(2), cart2.Items[0].Count)
  
  cart2 = c.DeductSKU(cart1.ID, cart1.Items[0].ID, 1)
  assert.Equal(t, float32(1), cart2.Items[0].Count)
  
  cart2 = c.DeductSKU(cart1.ID, cart1.Items[0].ID, 1)
  assert.Equal(t, int(0), len(cart2.Items))

  c.Close()
}

func TestLoadRedis(t *testing.T) {
  flag.Set("alsologtostderr", "true")
  flag.Set("log_dir", ".")
  flag.Set("v", "9")
  flag.Parse()

  glog.Info("Logging configured")
  c := New("redis://localhost:6379/0", 0, "", 100)
  res := c.Mode()
  assert.Equal(t, "redis", res)

  assert.Nil(t, c.GetCart("1"))
  cart1 := c.AppendSKU("", "product1", "", 1)
  
  cart2 := c.GetCart(cart1.ID)
  assert.NotNil(t, cart2)
  assert.Equal(t, int(1), len(cart2.Items))
  assert.Equal(t, float32(1), cart2.Items[0].Count)
  
  cart1 = c.AppendSKU(cart1.ID, "product1", "", 1)
  cart2 = c.GetCart(cart1.ID)
  assert.Equal(t, float32(2), cart2.Items[0].Count)
  
  cart2 = c.DeductSKU(cart1.ID, cart1.Items[0].ID, 1)
  assert.Equal(t, float32(1), cart2.Items[0].Count)
  
  cart2 = c.DeductSKU(cart1.ID, cart1.Items[0].ID, 1)
  assert.Equal(t, int(0), len(cart2.Items))

  c.Close()
}

