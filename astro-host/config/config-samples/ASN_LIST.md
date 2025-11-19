# China-Based ASN List

This file contains the list of Autonomous System Numbers (ASNs) used for fetching IPv6 routes.

## Major Carriers

### China Telecom (中国电信)
- **AS4134** - China Telecom (Primary)
- **AS4809** - China Telecom Next Carrying Network
- **AS23724** - China Telecom Group
- **AS58466** - China Telecom
- **AS58519** - China Telecom
- **AS58536** - China Telecom
- **AS58540** - China Telecom
- **AS58541** - China Telecom
- **AS58542** - China Telecom
- **AS58543** - China Telecom

### China Unicom (中国联通)
- **AS4837** - China Unicom (Primary)
- **AS9929** - China Unicom Industrial Internet Backbone
- **AS17621** - China Unicom
- **AS17622** - China Unicom
- **AS17623** - China Unicom

### China Mobile (中国移动)
- **AS9808** - China Mobile (Primary)
- **AS24400** - China Mobile Communications Corporation
- **AS56040** - China Mobile
- **AS56041** - China Mobile
- **AS56042** - China Mobile
- **AS56044** - China Mobile
- **AS56046** - China Mobile
- **AS56047** - China Mobile
- **AS56048** - China Mobile

## Education & Research Networks

### CERNET (中国教育和科研计算机网)
- **AS4538** - China Education and Research Network Center
- **AS23910** - CERNET2

### CSTNET (中国科技网)
- **AS7497** - Computer Network Information Center, Chinese Academy of Sciences

## Major Tech Companies

### Tencent (腾讯)
- **AS45090** - Tencent Building, Kejizhongyi Avenue
- **AS132203** - Tencent Building, Kejizhongyi Avenue

### Alibaba (阿里巴巴)
- **AS37963** - Alibaba (China) Technology Co., Ltd.
- **AS45102** - Alibaba (US) Technology Co., Ltd.

### Baidu (百度)
- **AS38365** - Beijing Baidu Netcom Science and Technology Co., Ltd.
- **AS55967** - Baidu International Technology (Shenzhen) Co., Ltd.

## Other Major ISPs & Networks

### Dr.Peng Telecom & Media Group (鹏博士)
- **AS17816** - China Unicom IP network China169 Guangdong province
- **AS17897** - Dr.Peng Telecom & Media Group Co., Ltd

### Regional & Specialized Networks
- **AS4229** - China Unicom (Not in Service)
- **AS4812** - China Telecom (Group)
- **AS4847** - China Networks Inter-Exchange
- **AS9394** - China TieTong Telecommunications Corporation
- **AS18118** - Shenzhen Topway Video Communication

### CDN & Cloud Providers
- **AS21859** - Zenlayer Inc
- **AS23764** - CTGNet
- **AS24407** - China Mobile Communications Corporation
- **AS24429** - Taobao (China) Software Co., Ltd.
- **AS24444** - Shanghaicloud
- **AS24482** - SG.GS

### Additional Networks
- **AS38019** - Super Micro Computer, Inc.
- **AS38283** - CHINANET-SZ-AP
- **AS38587** - China Mobile Hong Kong Company Limited
- **AS38854** - QingDao, ChengXinYunJiSuan LTD

### Smaller ISPs & Regional Providers
- **AS45062** - ShenZhen Topway Video Communication
- **AS45071** - GuangDong Mobile Communication
- **AS45075** - Shenzhen Topway Video Communication

### International & Overseas Networks
- **AS55720** - Gigabit Hosting Sdn Bhd
- **AS58453** - CMI Limited Mobile
- **AS58593** - Guangdong Mobile Communication Co.Ltd.
- **AS58772** - China Mobile Group Guangdong Company Limited
- **AS58834** - Shanghai Anchnet Network Technology
- **AS59019** - Hutchison Global Communications Limited

### Recent Allocations & Specialized Services
- **AS62610** - Zenlayer Limited
- **AS63525** - Shenzhen JingXunTong Technology Co., Ltd.
- **AS63593** - Eons Data Communications Limited
- **AS63629** - Shenzhen Tencent Computer Systems Company Limited
- **AS63655** - Beijing Volcano Engine Technology Co., Ltd.

### High-Range ASNs (132xxx-137xxx)
- **AS132525** - SoftLayer Dutch Holdings B.V.
- **AS133111** - CHINANET Tianjin province network
- **AS133219** - China Unicom Liaoning Province Network
- **AS133365** - China Unicom Shandong Province Network
- **AS133492** - China Telecom Guangxi province
- **AS134419** - CHINANET Fujian province network
- **AS134543** - CHINANET Sichuan province network
- **AS134756** - China Telecom Beijing Telecom
- **AS134770** - China Telecom Shanghai Telecom

### Additional High-Range ASNs
- **AS135377** - Ucloud Information Technology
- **AS136190** - China Unicom Hubei Province Network
- **AS136958** - China Mobile Group Jiangsu Company Limited
- **AS137266** - China Telecom Next Generation Carrier Network
- **AS137539** - Beijing Kingsoft Cloud Internet Technology
- **AS137687** - CHINANET Jiangsu province network
- **AS137710** - China Telecom Sichuan Telecom Internet Data Center
- **AS137718** - Beijing Jingdong 360 Du E-Commerce
- **AS137726** - China Mobile Group Guangdong Company Limited
- **AS137735** - Sichuan Mobile Communication Company

---

## Total ASNs: 80+

**Last Updated:** November 3, 2025

## Usage

These ASNs are used by:
- `fetch_routes.py` - Python-based route fetcher
- `fetch_fresh_routes.sh` - Bash-based route fetcher

To fetch routes for these ASNs:
```bash
./fetch_routes.py -v --aggregate
```

To use a custom ASN list:
```bash
./fetch_routes.py --asn 4134 --asn 4837 --asn 9808
```
