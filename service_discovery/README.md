# goDiscovery
Godiscovery is a simple service register&discovery backed by etcd.  This is based on the code of [bytemode](https://blog.csdn.net/blogsun/article/details/102861648).
In addition, `consistent` is a consistent hash circle backed by treemap(modified on `stathat.com/c/consistent`)

## How It Works

1. register:
   * get a lease(default ttl 3s) from etcd
   * put the key(default  prefix+hostname) with the lease
   * keep alive the lease repeatedly
   * if service aborts, the lease will be expired. And then, the key will be removed

2. discovery:
   * get all keys under the same prefix, add them to hash circle
   * watching these keys. if a new key is put, add it to circle; if a key is removed,  delete it from circle

3. circle(modified on `stathat.com/c/consistent`):
    * backed by treemap (https://github.com/emirpasic/gods)
    * support murmur3 hash (https://github.com/spaolacci/murmur3)
    
    

## Usage Example

```go
func main() {
	_ = g.DefaultServiceRegister(endpoint, User, Password, Prefix)
	d := g.DefaultServiceDiscovery(endpoint, User, Password, Prefix, replicaNum, hash)

	node, err :=d.Allocate("user1")
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(node)
}
```





