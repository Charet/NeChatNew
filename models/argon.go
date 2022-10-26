package models

type CryptoParam struct {
	Memory      uint32 //算法使用的内存量(单位:KB)
	Iterations  uint32 //在内存上的迭代（或传递）次数
	Parallelism uint8  //算法使用的线程（或通道）的数量
	SaltLength  uint32 //随机盐的长度 建议密码散列使用16字节
	KeyLength   uint32 //生成的密钥（或密码散列）的长度 建议16字节或更多
}
