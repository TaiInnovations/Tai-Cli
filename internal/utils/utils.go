package utils

import "encoding/json"

func DeepCopy(src, dst interface{}) error {
    data, err := json.Marshal(src) // 序列化为 JSON
    if err != nil {
        return err
    }
    return json.Unmarshal(data, dst) // 反序列化为新结构
}
