package config

import "github.com/net12labs/cirm/ops/rtm"

func Init() {

	rtm.Etc.SetKV("unit_id", "default")
	rtm.Etc.SetKV("rtm_name", "china-ip-routes-maker")
	rtm.Etc.SetKV("domain_name", rtm.Etc.GetJoined("/", "rtm_name", "unit_id"))
	rtm.Etc.SetKV("home_dir", "../units/"+rtm.Etc.Get("unit_id").String())
	rtm.Etc.SetKV("pid_file_path", rtm.Etc.Get("home_dir").String()+"/proc/china-ip-routes-maker.pid")
	rtm.Etc.SetKV("data_dir", rtm.Etc.Get("home_dir").String()+"/data")
	rtm.Etc.SetKV("dom_db_path", rtm.Etc.Get("data_dir").String()+"/dom.db")
	rtm.Etc.SetKV("host_db_path", rtm.Etc.Get("data_dir").String()+"/host.db")
	rtm.Etc.SetKV("site_db_path", rtm.Etc.Get("data_dir").String()+"/site.db")
	rtm.Etc.SetKV("socket_path", rtm.Etc.Get("home_dir").String()+"/main.sock")
}
