package models

// Series acts like a cross between a `map[int]interface{}` and an `[]interface{}`.
type Series struct {
	Key  int           `bson:"key"  json:"key"`
	Data []interface{} `bson:"data" json:"data"`
}

func (self *Series) Set(key int, val interface{}) {
	if self.Data == nil {
		self.Data = []interface{}{val}
		self.Key = key
		return
	}

	maxKey := self.Key + len(self.Data)
	data := self.Data

	switch true {
	case key < self.Key:
		data = make([]interface{}, maxKey-key)
		copy(data[self.Key-key:maxKey-key], self.Data)
		self.Key = key
		self.Data = data
	case key >= maxKey:
		maxKey = key + 1
		data = make([]interface{}, maxKey-self.Key)
		copy(data, self.Data)
		self.Data = data
	}
	data[key-self.Key] = val
}

func (self *Series) Get(key int) interface{} {
	if self.Data == nil {
		return nil
	}
	i := key - self.Key
	if i >= 0 && i < len(self.Data) {
		return self.Data[i]
	}
	return nil
}

func (self *Series) Each(fn func(int, interface{})) {
	if self.Data == nil {
		return
	}
	for i, data := range self.Data {
		if data != nil {
			fn(self.Key+i, data)
		}
	}
}
