package client

import (
  "context"
  "fmt"
  "strings"
  "sync"
  "crypto/rand"

  "github.com/libp2p/go-libp2p"
  "github.com/libp2p/go-libp2p/core/crypto"
  "github.com/libp2p/go-libp2p/core/host"
  "github.com/libp2p/go-libp2p/core/peer"
  pubsub "github.com/libp2p/go-libp2p-pubsub"
  "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
  "github.com/libp2p/go-libp2p/core/protocol"
  "github.com/libp2p/go-libp2p/p2p/discovery/routing"
  "github.com/multiformats/go-multiaddr"
  "github.com/libp2p/go-libp2p-kad-dht"

  "github.com/golang/glog"
)

type P2PConnect struct {
  ctx          context.Context
  host         host.Host
	//rpcServer   *rpc.Server
	//rpcClient   *rpc.Client
  protocol     protocol.ID
  dhtRD        *routing.RoutingDiscovery
  gossipSub    *pubsub.PubSub 
  pubTopics     map[string]*pubsub.Topic
  subTopics     map[string]*pubsub.Topic
  
  mDNSService
  mDNSServer
}

type addrList []multiaddr.Multiaddr

func NewAddrList() *addrList {
  return &addrList{}
}

func (al *addrList) Len() int {
  return len(*al)
}

func (al *addrList) String() string {
	strs := make([]string, len(*al))
	for i, addr := range *al {
		strs[i] = addr.String()
	}
	return strings.Join(strs, ",")
}

func (al *addrList) Set(value string) error {
	addr, err := multiaddr.NewMultiaddr(value)
	if err != nil {
    glog.Errorf("ERR: error NewMultiaddr %s: %v", value, err)
		return err
	}
	*al = append(*al, addr)
	return nil
}

// discoveryNotifee gets notified when we find a new peer via mDNS discovery
type discoveryNotifee struct {
    h host.Host
}

// HandlePeerFound connects to peers discovered via mDNS. Once they're connected,
// the PubSub system will automatically start interacting with them if they also
// support PubSub.
func (n *discoveryNotifee) HandlePeerFound(pi peer.AddrInfo) {
  glog.Infof("LOG: discovered new peer %s", pi.ID.Pretty())
  err := n.h.Connect(context.Background(), pi)
  if err != nil {
    glog.Errorf("ERR: error connecting to peer %s: %s", pi.ID.Pretty(), err)
  }
}

func NewP2PConnect() *P2PConnect {
  p2p := &P2PConnect{
                     ctx: context.Background(),
                     pubTopics: make(map[string]*pubsub.Topic),
                     subTopics: make(map[string]*pubsub.Topic),
                  }
  //p2p.ctx, p2p.cancel = context.WithCancel(context.Background())
  //p2p.ctx = context.Background()
  return p2p
}

func (p2p *P2PConnect) Connected() bool {
  return p2p != nil && p2p.dhtRD != nil
}

func (p2p *P2PConnect) GetHostID() string {
  return p2p.host.ID().Pretty()
}

func newAddrsFactory(advertiseAddrs *addrList) func([]multiaddr.Multiaddr) []multiaddr.Multiaddr {
	return func([]multiaddr.Multiaddr) []multiaddr.Multiaddr {
		return []multiaddr.Multiaddr(*advertiseAddrs)
	}
}

func (p2p *P2PConnect) NewHost(port int, bootstrapPeers *addrList, privKey crypto.PrivKey) bool {
  pk := privKey
  if privKey == nil {
    priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
    if err != nil {
      glog.Errorf("ERR: GenerateKeyPairWithReader: %v", err)
      return false
    }
    pk = priv
  }

  var err error
	addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
  
	p2p.host, err = libp2p.New(//p2p.ctx,
                  libp2p.ListenAddrs(addr),
                  libp2p.Identity(pk),
                  libp2p.NATPortMap(),
                  libp2p.AddrsFactory(newAddrsFactory(bootstrapPeers)),
                )
	if err != nil {
    glog.Errorf("ERR: libp2p.New: %v", err)
	}
  return err == nil
}


// setupDiscovery creates an mDNS discovery service and attaches it to the libp2p Host.
// This lets us automatically discover peers on the same LAN and connect to them.
func (p2p *P2PConnect) NewMDNS(discoveryServiceTag string, commonName string) bool {
  // setup mDNS discovery to find local peers
  glog.Infof("INFO: NewMdnsService: %s", discoveryServiceTag + commonName)
  p2p.mDNSService = mdns.NewMdnsService(p2p.host, discoveryServiceTag + commonName, &discoveryNotifee{h: p2p.host})
  err := p2p.mDNSService.Start()
	if err != nil {
    glog.Errorf("ERR: libp2p.NewMdnsService: %v", err)
	}
  return err == nil
}


func (p2p *P2PConnect) NewKDHT(bootstrapPeers *addrList) bool {
	var options []dht.Option
  glog.Infof("LOG: NewKDHT with bootstrap node: %v", *bootstrapPeers)
	if bootstrapPeers.Len() == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
	}

	kdht, err := dht.New(p2p.ctx, p2p.host, options...)
	if err != nil {
    glog.Errorf("ERR: dht.New: %v", err)
		return false
	}

	if err = kdht.Bootstrap(p2p.ctx); err != nil {
    glog.Errorf("ERR: dht.Bootstrap: %v", err)
		return false
	}
	var wg sync.WaitGroup
	for _, peerAddr := range (*bootstrapPeers) {
		peerinfo, errp := peer.AddrInfoFromP2pAddr(peerAddr)
    if errp != nil {
      glog.Errorf("ERR: peer.AddrInfoFromP2pAddr: %v", errp)
      continue
    }
    glog.Infof("LOG: Connecting with bootstrap node: %v", peerinfo)
    
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := p2p.host.Connect(p2p.ctx, *peerinfo); err != nil {
				glog.Errorf("ERR: Error while connecting to node %q: %-v", peerinfo, err)
			} else {
				glog.Infof("LOG: Connection established with bootstrap node: %q", *peerinfo)
			}
		}()
	}
	wg.Wait()
  glog.Infof("INFO: RoutingDiscovery")
  p2p.dhtRD = routing.NewRoutingDiscovery(kdht)
	return true
}

func (p2p *P2PConnect) NewGossipSub() bool {
  // create a new PubSub service using the GossipSub router
  var err error
  glog.Infof("INFO: New GossipSub")
  p2p.gossipSub, err = pubsub.NewGossipSub(p2p.ctx, p2p.host)
  if err != nil {
    glog.Errorf("ERR: NewGossipSub: %v", err)
    return false
  }
  return true
}

// start subsriber to topic
func (p2p *P2PConnect) Subscribe(topicName string, event func(*pubsub.Message)) bool {
  glog.Infof("INFO: Join Topic: %s", topicName)
  topic, err := p2p.gossipSub.Join(topicName)
  if err != nil {
    glog.Errorf("ERR: gossipSub.Join: %v", err)
    return false
  }

  // subscribe to topic
  subscriber, err := topic.Subscribe()
  if err != nil {
    glog.Errorf("ERR: topic.Subscribe: %v", err)
    return false
  }
  glog.Infof("INFO: Listen Topic: %s", topicName)
  for {
    msg, err := subscriber.Next(p2p.ctx)
    glog.Infof("INFO: Listen Topic: %s: %v: %v", topicName, msg, err)
    if err != nil {
      glog.Errorf("ERR: subscriber.Next: %v", err)
      return false
    }
    glog.Infof("LOG: Got message: %s, from: %s", string(msg.Data), msg.ReceivedFrom.Pretty())
    // only consider messages delivered by other peers
    if msg.ReceivedFrom == p2p.host.ID() {
      continue
    }
    event(msg)
    //fmt.Printf("[2] got message: %s, from: %s\n", string(msg.Data), msg.ReceivedFrom.Pretty())
  }
  glog.Infof("INFO: EXIT Listen Topic: %s", topicName)
  return true
}

func (p2p *P2PConnect) Publish(topicName string, msg []byte) bool {
  var err error
  topic, ok := p2p.pubTopics[topicName]
  if !ok {
    topic, err = p2p.gossipSub.Join(topicName)
    glog.Infof("INFO: Join Topic: %s", topicName)
    if err != nil {
      glog.Errorf("ERR: Join: %v", err)
      return false
    }
    p2p.pubTopics[topicName] = topic
  }
  err = topic.Publish(p2p.ctx, msg)
  if err != nil {
		glog.Errorf("ERR: Publish(%s): %v", topicName, err)
    return false
	}
  glog.Infof("LOG: Publish(%s): %v", topicName, err)
  return true
}
