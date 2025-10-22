package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// getLocalIP retrieves the local IP address and subnet mask for the WiFi interface
func getLocalIP() (net.IP, *net.IPNet, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, ipnet, nil
			}
		}
	}
	return nil, nil, fmt.Errorf("no suitable IP address found")
}

// scanIP checks if an IP address is active by attempting to connect to common ports
func scanIP(ip string, wg *sync.WaitGroup) {
	defer wg.Done()
	ports := []int{22, 80, 443} // Common ports: SSH, HTTP, HTTPS
	for _, port := range ports {
		address := fmt.Sprintf("%s:%d", ip, port)
		conn, err := net.DialTimeout("tcp", address, time.Second*1)
		if err == nil {
			conn.Close()
			fmt.Printf("Active IP: %s (port %d open)\n", ip, port)
			return // If any port is open, consider active and stop checking
		}
	}
}

func main() {
	localIP, ipNet, err := getLocalIP()
	if err != nil {
		fmt.Printf("Error getting local IP: %v\n", err)
		return
	}

	fmt.Printf("Local IP: %s, Subnet: %s\n", localIP, ipNet)

	// Calculate IP range (assuming /24 for simplicity)
	ipRange := make([]string, 0)
	for i := 1; i < 255; i++ {
		ip := net.IP{localIP[0], localIP[1], localIP[2], byte(i)}
		ipRange = append(ipRange, ip.String())
	}

	var wg sync.WaitGroup
	for _, ip := range ipRange {
		wg.Add(1)
		go scanIP(ip, &wg)
	}
	wg.Wait()
	fmt.Println("Scan complete.")
}
