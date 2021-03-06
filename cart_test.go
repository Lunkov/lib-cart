package cart

import (
  "testing"
  "github.com/stretchr/testify/assert"

  "fmt"
  "strconv"
  "math/rand"
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
  assert.Equal(t, cart1.ID, cart2.ID)
  assert.Equal(t, float32(2), cart2.Items[0].Count)
  
  cart2 = c.DeductSKU(cart1.ID, cart1.Items[0].ID, 1)
  assert.Equal(t, cart1.ID, cart2.ID)
  assert.Equal(t, float32(1), cart2.Items[0].Count)
  
  cart2 = c.DeductSKU(cart1.ID, cart1.Items[0].ID, 1)
  // assert.Equal(t, cart1.ID, cart2.ID)
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
  /* TODO
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
  assert.Equal(t, cart1.ID, cart2.ID)
  assert.Equal(t, float32(2), cart2.Items[0].Count)
  
  cart2 = c.DeductSKU(cart1.ID, cart1.Items[0].ID, 1)
  assert.Equal(t, float32(1), cart2.Items[0].Count)
  
  cart2 = c.DeductSKU(cart1.ID, cart1.Items[0].ID, 1)
  assert.Equal(t, int(0), len(cart2.Items))
  */
  c.Close()
}


func BenchmarkCartMemomry(b *testing.B) {
  flag.Set("alsologtostderr", "true")
  flag.Set("log_dir", ".")
  flag.Set("v", "0")
  flag.Parse()
  
  cartCount := int64(1000)
  productCount := int64(100)
  c := New("memory", 1000, "", 100)
  res := c.Mode()
  assert.Equal(b, "memory", res)
  
  b.ResetTimer()
  for i := 1; i <= 8; i *= 2 {
		b.Run(strconv.Itoa(i), func(b *testing.B) {
			b.SetParallelism(i)
      b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
          cartID := fmt.Sprintf("CART-%d", rand.Int63n(cartCount))
          sku := fmt.Sprintf("SKU-%d", rand.Int63n(productCount))
          cart2 := c.AppendSKU(cartID, sku, "", 1)
          
          assert.Equal(b, cartID, cart2.ID)
        }
      })
    })
  }
  c.Close()
}
