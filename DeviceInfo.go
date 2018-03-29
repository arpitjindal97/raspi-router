package main

func DeviceInfo() OSInfo {
	var info OSInfo

	go func() {
		info.DistributionId = GetOutput("lsb_release -i | awk '{print $3}'")
		info.Description = GetOutput("lsb_release -d | cut -f1 --complement'")
		info.Release = GetOutput(" lsb_release -r | cut -f2")
	}()

	info.Codename = GetOutput(" lsb_release -c | cut -f2")

	go func() {
		info.Hostname = GetOutput("hostname")
		info.KernelRelease = GetOutput("uname -r")
		info.Architecture = GetOutput("uname -m")
		info.ModelName = GetOutput("lscpu | grep 'Model name' | awk '{$1=\"\";$2=\"\";print}'")
		info.ModelName = info.ModelName[2:]
	}()

	info.CPUs = GetOutput("lscpu | grep CPU\\(s\\): | grep -v node | awk '{print $2}'")

	info.LocalTime = GetOutput("date")
	info.TimeZone = GetOutput("timedatectl | grep zone | awk '{$1=\"\";$2=\"\";print}'")
	info.TimeZone = info.TimeZone[2:]

	info.UpTime = GetOutput("uptime -p")
	info.UpSince = GetOutput("uptime -s")
	return info
}
