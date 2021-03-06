package cart

import (
  "github.com/golang/glog"
  "github.com/google/uuid"
  
  "github.com/Lunkov/lib-cache"
)

type Cart struct {
  cache *cache.Cache
}

type ItemsInfo struct {
  ID         string
  ParentID   string
  SKU        string
  Count      float32  
}

type CartInfo struct {
  ID       string
  Items  []ItemsInfo
}

func genUUID() string {
  uid, _ := uuid.NewUUID()
  return uid.String()
}

func (c *Cart) HasError() bool {
  if c.cache == nil {
    return true
  }
  return c.cache.HasError()
}

func (c *Cart) Mode() string {
  if c.cache == nil {
    return "undefined"
  }
  return c.cache.Mode()
}

func (c *Cart) Count() int64 {
  if c.cache == nil {
    return -1
  }
  return c.cache.Count()
}

func (c *Cart) DestroyAll() {
  if c.cache != nil {
    c.cache.Clear()
  }
}

func (c *Cart) GetCart(id string) *CartInfo {
  var ci CartInfo
  ii, ok := c.cache.Get(id, &ci)
  if !ok {
    return nil
  }
  ci, ok = (ii).(CartInfo)
  if !ok {
    return nil
  }
  return &ci
}

func (c *Cart) AppendSKU(cartId string, sku string, parentId string, count float32) *CartInfo {
  var ci CartInfo
  ii, ok := c.cache.Get(cartId, &ci)
  f := false
  if ok {
    ci, ok = (ii).(CartInfo)
    for k, v := range ci.Items {
      if v.SKU == sku && v.ParentID == parentId {
        v.Count += count
        ci.Items[k] = v
        f = true
        break
      }
    }
  }
  if !ok {
    if cartId == "" {
      cartId = genUUID()
    }
    ci.ID = cartId
    ci.Items = make([]ItemsInfo, 0, 0)
  }
  if !f {
    ci.Items = append(ci.Items, ItemsInfo{ID: genUUID(), ParentID: parentId, SKU: sku, Count: count})
  }
  c.cache.Set(cartId, ci)
  return &ci
}

func remove(s []ItemsInfo, i int) []ItemsInfo {
  s[len(s)-1], s[i] = s[i], s[len(s)-1]
  return s[:len(s)-1]
}

func (c *Cart) DeductSKU(cartId string, itemId string, count float32) *CartInfo {
  var ci CartInfo
  ii, ok := c.cache.Get(cartId, &ci)
  if ok {
    ci, ok = (ii).(CartInfo)
    if ok {
      for k, v := range ci.Items {
        if v.ID == itemId {
          v.Count -= count
          if v.Count > 0 {
            ci.Items[k] = v
          } else {
            remove(ci.Items, k)
          }
          c.cache.Set(cartId, &ci)
          break
        }
      }
    }
  } else {
    if cartId == "" {
      cartId = genUUID()
    }
    ci.ID = cartId
    ci.Items = make([]ItemsInfo, 0, 0)
    c.cache.Set(cartId, &ci)
  }
  return &ci
}

////
// Init
////
func New(mode string, expiryTime int64, URL string, MaxConnections int) *Cart {
  if glog.V(9) {
    glog.Infof("DBG: CART: Init")
  }
  c := &Cart{}
  c.cache = cache.New(mode, expiryTime, URL, MaxConnections)
  if c.cache == nil {
    glog.Errorf("ERR: CART: Init(%s) error", mode)
    return nil
  }
  glog.Infof("LOG: CART: Mode is %s", c.cache.Mode())
  return c
}

func (c *Cart) Close() {
  if c.cache != nil {
    c.cache.Close()
    c.cache = nil
  }
}
