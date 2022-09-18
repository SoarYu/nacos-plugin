/*
 * Copyright 1999-2018 Alibaba Group Holding Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *      http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nacos

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"net"
)

type Nacos struct {
	Next  plugin.Handler
	Zones []string
}

func (vs *Nacos) String() string {
	b, err := json.Marshal(vs)

	if err != nil {
		return ""
	}

	return string(b)
}

func (vs *Nacos) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	state := request.Request{W: w, Req: r}

	name := state.QName()

	m := new(dns.Msg)

	clientIP := state.IP()
	if clientIP == "127.0.0.1" {
		clientIP = LocalIP()
	}

	if host, err := GrpcClient.SelectOneHealthyInstances(name[:len(name)-1]); host == nil || err != nil {
		return plugin.NextOrFailure(vs.Name(), vs.Next, ctx, w, r)
	} else {
		answer := make([]dns.RR, 0)
		extra := make([]dns.RR, 0)
		var rr dns.RR

		switch state.Family() {
		case 1:
			rr = new(dns.A)
			rr.(*dns.A).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeA, Class: state.QClass(), Ttl: DNSTTL}
			rr.(*dns.A).A = net.ParseIP(host.Ip).To4()
		case 2:
			rr = new(dns.AAAA)
			rr.(*dns.AAAA).Hdr = dns.RR_Header{Name: state.QName(), Rrtype: dns.TypeAAAA, Class: state.QClass(), Ttl: DNSTTL}
			rr.(*dns.AAAA).AAAA = net.ParseIP(host.Ip)
		}

		srv := new(dns.SRV)
		srv.Hdr = dns.RR_Header{Name: "_" + state.Proto() + "." + state.QName(), Rrtype: dns.TypeSRV, Class: state.QClass(), Ttl: DNSTTL}
		port := host.Port
		srv.Port = uint16(port)
		srv.Target = "."

		extra = append(extra, srv)
		answer = append(answer, rr)

		m.Answer = answer
		m.Extra = extra
		result, _ := json.Marshal(m.Answer)
		//NacosClientLogger.Info("[RESOLVE]", " ["+name[:len(name)-1]+"]  result: "+string(result)+", clientIP: "+clientIP)
		fmt.Println("[RESOLVE]", " ["+name[:len(name)-1]+"]  result: "+string(result)+", clientIP: "+clientIP)
	}

	m.SetReply(r)
	m.Authoritative, m.RecursionAvailable, m.Compress = true, true, true

	state.SizeAndDo(m)
	m = state.Scrub(m)
	w.WriteMsg(m)
	return dns.RcodeSuccess, nil
}

func (vs *Nacos) Name() string { return "nacos" }
