package tconfig
//
//import (
//	"bytes"
//	"context"
//	"github.com/coreos/etcd/clientv3"
//	"github.com/spf13/viper"
//	"io"
//	"time"
//)
//
//// 跳过 crypt，自己实现加载远程配置的功能，其实就是实现 viper 中的 remoteConfigFactory 接口：
//// https://github.com/spf13/viper/blob/master/remote/remote.go
//type Config struct {
//	viper.RemoteProvider
//
//	Username string
//	Password string
//}
//
//func (c *Config) Get(rp viper.RemoteProvider) (io.Reader, error) {
//	c.RemoteProvider = rp
//
//	return c.get()
//}
//
//func (c *Config) Watch(rp viper.RemoteProvider) (io.Reader, error) {
//	c.RemoteProvider = rp
//
//	return c.get()
//}
//
//func (c *Config) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
//	c.RemoteProvider = rp
//
//	rr := make(chan *viper.RemoteResponse)
//	stop := make(chan bool)
//
//	go func() {
//		for {
//			client, err := c.newClient()
//
//			if err != nil {
//				time.Sleep(time.Duration(time.Second))
//				continue
//			}
//
//			defer client.Close()
//
//			ch := client.Watch(context.Background(), c.RemoteProvider.Path())
//
//			select {
//			case <-stop:
//				return
//			case res := <-ch:
//				for _, event := range res.Events {
//					rr <- &viper.RemoteResponse{
//						Value: event.Kv.Value,
//					}
//				}
//			}
//		}
//	}()
//
//	return rr, stop
//}
//
//func (c *Config) newClient() (*clientv3.Client, error) {
//	client, err := clientv3.New(clientv3.Config{
//		Endpoints: []string{c.Endpoint()},
//		Username:  c.Username,
//		Password:  c.Password,
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	return client, nil
//}
//
//func (c *Config) get() (io.Reader, error) {
//	client, err := c.newClient()
//
//	if err != nil {
//		return nil, err
//	}
//
//	defer client.Close()
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	resp, err := client.Get(ctx, c.Path())
//	cancel()
//
//	if err != nil {
//		return nil, err
//	}
//	if resp.Kvs == nil {
//		return bytes.NewReader(make([]byte, 0)), nil
//	}
//	return bytes.NewReader(resp.Kvs[0].Value), nil
//}
