package client

import "encoding/json"

type PingResponse struct {
	BuildVersion string `json:"build_version"`
	Version      string `json:"version"`
	PID          int    `json:"pid"`
}

/*

{
  "config": {
    "cpus": {
      "boot_vcpus": 1,
      "max_vcpus": 1,
      "topology": {
        "threads_per_core": 0,
        "cores_per_die": 0,
        "dies_per_package": 0,
        "packages": 0
      },
      "kvm_hyperv": false,
      "max_phys_bits": 0,
      "affinity": [
        {
          "vcpu": 0,
          "host_cpus": [
            0
          ]
        }
      ],
      "features": {
        "amx": true
      }
    },
    "memory": {
      "size": "512 MB",
      "hotplug_size": 0,
      "hotplugged_size": 0,
      "mergeable": false,
      "hotplug_method": "Acpi",
      "shared": false,
      "hugepages": false,
      "hugepage_size": 0,
      "prefault": false,
      "thp": true,
      "zones": [
        {
          "id": "string",
          "size": "512 MB",
          "file": "string",
          "mergeable": false,
          "shared": false,
          "hugepages": false,
          "hugepage_size": 0,
          "host_numa_node": 0,
          "hotplug_size": 0,
          "hotplugged_size": 0,
          "prefault": false
        }
      ]
    },
    "payload": {
      "firmware": "string",
      "kernel": "string",
      "cmdline": "string",
      "initramfs": "string"
    },
    "disks": [
      {
        "path": "string",
        "readonly": false,
        "direct": false,
        "iommu": false,
        "num_queues": 1,
        "queue_size": 128,
        "vhost_user": false,
        "vhost_socket": "string",
        "rate_limiter_config": {
          "bandwidth": {
            "size": 0,
            "one_time_burst": 0,
            "refill_time": 0
          },
          "ops": {
            "size": 0,
            "one_time_burst": 0,
            "refill_time": 0
          }
        },
        "pci_segment": 0,
        "id": "string"
      }
    ],
    "net": [
      {
        "tap": "string",
        "ip": "192.168.249.1",
        "mask": "255.255.255.0",
        "mac": "string",
        "host_mac": "string",
        "mtu": 0,
        "iommu": false,
        "num_queues": 2,
        "queue_size": 256,
        "vhost_user": false,
        "vhost_socket": "string",
        "vhost_mode": "Client",
        "id": "string",
        "pci_segment": 0,
        "rate_limiter_config": {
          "bandwidth": {
            "size": 0,
            "one_time_burst": 0,
            "refill_time": 0
          },
          "ops": {
            "size": 0,
            "one_time_burst": 0,
            "refill_time": 0
          }
        }
      }
    ],
    "rng": {
      "src": "/dev/urandom",
      "iommu": false
    },
    "balloon": {
      "size": 0,
      "deflate_on_oom": false,
      "free_page_reporting": false
    },
    "fs": [
      {
        "tag": "string",
        "socket": "string",
        "num_queues": 1,
        "queue_size": 1024,
        "pci_segment": 0,
        "id": "string"
      }
    ],
    "pmem": [
      {
        "file": "string",
        "size": 0,
        "iommu": false,
        "discard_writes": false,
        "pci_segment": 0,
        "id": "string"
      }
    ],
    "serial": {
      "file": "string",
      "mode": "Off",
      "iommu": false
    },
    "console": {
      "file": "string",
      "mode": "Off",
      "iommu": false
    },
    "devices": [
      {
        "path": "string",
        "iommu": false,
        "pci_segment": 0,
        "id": "string"
      }
    ],
    "vdpa": [
      {
        "path": "string",
        "num_queues": 1,
        "iommu": false,
        "pci_segment": 0,
        "id": "string"
      }
    ],
    "vsock": {
      "cid": 3,
      "socket": "string",
      "iommu": false,
      "pci_segment": 0,
      "id": "string"
    },
    "sgx_epc": [
      {
        "id": "string",
        "size": 0,
        "prefault": false
      }
    ],
    "numa": [
      {
        "guest_numa_id": 0,
        "cpus": [
          0
        ],
        "distances": [
          {
            "destination": 0,
            "distance": 0
          }
        ],
        "memory_zones": [
          "string"
        ],
        "sgx_epc_sections": [
          "string"
        ]
      }
    ],
    "iommu": false,
    "watchdog": false,
    "platform": {
      "num_pci_segments": 0,
      "iommu_segments": [
        0
      ],
      "serial_number": "string",
      "uuid": "string",
      "oem_strings": [
        "string"
      ],
      "tdx": false
    },
    "tpm": {
      "socket": "string"
    }
  },
  "state": "Created",
  "memory_actual_size": 0,
  "device_tree": {
    "additionalProp1": {
      "id": "string",
      "resources": [
        {}
      ],
      "children": [
        "string"
      ],
      "pci_bdf": "string"
    },
    "additionalProp2": {
      "id": "string",
      "resources": [
        {}
      ],
      "children": [
        "string"
      ],
      "pci_bdf": "string"
    },
    "additionalProp3": {
      "id": "string",
      "resources": [
        {}
      ],
      "children": [
        "string"
      ],
      "pci_bdf": "string"
    }
  }
}
*/

type VMConfig struct {
	CPUs     *VMConfigCPUs     `json:"cpus,omitempty"`
	Memory   *VMConfigMemory   `json:"memory,omitempty"`
	Payload  *VMConfigPayload  `json:"payload,omitempty"`
	Disks    []*VMConfigDisk   `json:"disks,omitempty"`
	Net      []*VMConfigNet    `json:"net,omitempty"`
	RNG      *VMConfigRNG      `json:"rng,omitempty"`
	Balloon  *VMConfigBalloon  `json:"balloon,omitempty"`
	FS       []*VMConfigFS     `json:"fs,omitempty"`
	PMEM     []*VMConfigPMEM   `json:"pmem,omitempty"`
	Serial   *VMConfigSerial   `json:"serial,omitempty"`
	Console  *VMConfigConsole  `json:"console,omitempty"`
	Devices  []*VMConfigDevice `json:"devices,omitempty"`
	VDPA     []*VMConfigVDPA   `json:"vdpa,omitempty"`
	VSOCK    *VMConfigVSOCK    `json:"vsock,omitempty"`
	SGXEPC   []*VMConfigSGXEPC `json:"sgx_epc,omitempty"`
	Numa     []*VMConfigNuma   `json:"numa,omitempty"`
	IOMMU    bool              `json:"iommu,omitempty"`
	Watchdog bool              `json:"watchdog,omitempty"`
	Platform *VMConfigPlatform `json:"platform,omitempty"`
	TPM      *VMConfigTPM      `json:"tpm,omitempty"`
}

type VMConfigTPM struct {
	Socket string `json:"socket,omitempty"`
}

type VMConfigPlatform struct {
	NumPCISegments int      `json:"num_pci_segments,omitempty"`
	IOMMUSegments  []int    `json:"iommu_segments,omitempty"`
	SerialNumber   string   `json:"serial_number,omitempty"`
	UUID           string   `json:"uuid,omitempty"`
	OEMStrings     []string `json:"oem_strings,omitempty"`
	TDX            bool     `json:"tdx,omitempty"`
}

type VMConfigNuma struct {
	GuestNumaID    int                     `json:"guest_numa_id,omitempty"`
	CPUs           []int                   `json:"cpus,omitempty"`
	Distances      []*VMConfigNumaDistance `json:"distances,omitempty"`
	MemoryZones    []string                `json:"memory_zones,omitempty"`
	SGXEPCSections []string                `json:"sgx_epc_sections,omitempty"`
}

type VMConfigNumaDistance struct {
	Destination int `json:"destination,omitempty"`
	Distance    int `json:"distance,omitempty"`
}

type VMConfigSGXEPC struct {
	ID       string `json:"id,omitempty"`
	Size     int    `json:"size,omitempty"`
	Prefault bool   `json:"prefault,omitempty"`
}

type VMConfigVSOCK struct {
	CID        int    `json:"cid,omitempty"`
	Socket     string `json:"socket,omitempty"`
	IOMMU      bool   `json:"iommu,omitempty"`
	PCISegment int    `json:"pci_segment,omitempty"`
	ID         string `json:"id,omitempty"`
}

type VMConfigVDPA struct {
	Path       string `json:"path,omitempty"`
	NumQueues  int    `json:"num_queues,omitempty"`
	IOMMU      bool   `json:"iommu,omitempty"`
	PCISegment int    `json:"pci_segment,omitempty"`
	ID         string `json:"id,omitempty"`
}

type VMConfigDevice struct {
	Path       string `json:"path,omitempty"`
	IOMMU      bool   `json:"iommu,omitempty"`
	PCISegment int    `json:"pci_segment,omitempty"`
	ID         string `json:"id,omitempty"`
}

type VMConfigConsole struct {
	File  string `json:"file,omitempty"`
	Mode  string `json:"mode,omitempty"`
	IOMMU bool   `json:"iommu,omitempty"`
}

type VMConfigSerial struct {
	File  string `json:"file,omitempty"`
	Mode  string `json:"mode,omitempty"`
	IOMMU bool   `json:"iommu,omitempty"`
}

type VMConfigPMEM struct {
	File          string `json:"file,omitempty"`
	Size          int    `json:"size,omitempty"`
	IOMMU         bool   `json:"iommu,omitempty"`
	DiscardWrites bool   `json:"discard_writes,omitempty"`
	PCISegment    int    `json:"pci_segment,omitempty"`
	ID            string `json:"id,omitempty"`
}

type VMConfigFS struct {
	Tag        string `json:"tag,omitempty"`
	Socket     string `json:"socket,omitempty"`
	NumQueues  int    `json:"num_queues,omitempty"`
	QueueSize  int    `json:"queue_size,omitempty"`
	PCISegment int    `json:"pci_segment,omitempty"`
	ID         string `json:"id,omitempty"`
}

type VMConfigBalloon struct {
	Size              int  `json:"size,omitempty"`
	DeflateOnOOM      bool `json:"deflate_on_oom,omitempty"`
	FreePageReporting bool `json:"free_page_reporting,omitempty"`
}

type VMConfigRNG struct {
	SRC   string `json:"src,omitempty"`
	IOMMU bool   `json:"iommu,omitempty"`
}

type VMConfigNet struct {
	Tap               string                     `json:"tap,omitempty"`
	IP                string                     `json:"ip,omitempty"`
	Mask              string                     `json:"mask,omitempty"`
	MAC               string                     `json:"mac,omitempty"`
	HostMAC           string                     `json:"host_mac,omitempty"`
	MTU               int                        `json:"mtu,omitempty"`
	IOMMU             bool                       `json:"iommu,omitempty"`
	NumQueues         int                        `json:"num_queues,omitempty"`
	QueueSize         int                        `json:"queue_size,omitempty"`
	VhostUser         bool                       `json:"vhost_user,omitempty"`
	VhostSocket       string                     `json:"vhost_socket,omitempty"`
	VhostMode         string                     `json:"vhost_mode,omitempty"`
	ID                string                     `json:"id,omitempty"`
	PCISegment        int                        `json:"pci_segment,omitempty"`
	RateLimiterConfig *VMConfigRateLimiterConfig `json:"rate_limiter_config,omitempty"`
}

type VMConfigDisk struct {
	Path              string                     `json:"path,omitempty"`
	Readonly          bool                       `json:"readonly,omitempty"`
	Direct            bool                       `json:"direct,omitempty"`
	IOMMU             bool                       `json:"iommu,omitempty"`
	NumQueues         int                        `json:"num_queues,omitempty"`
	QueueSize         int                        `json:"queue_size,omitempty"`
	VhostUser         bool                       `json:"vhost_user,omitempty"`
	VhostSocket       string                     `json:"vhost_socket,omitempty"`
	RateLimiterConfig *VMConfigRateLimiterConfig `json:"rate_limiter_config,omitempty"`
	PCISegment        int                        `json:"pci_segment,omitempty"`
	ID                string                     `json:"id,omitempty"`
}

type VMConfigRateLimiterConfig struct {
	Bandwidth *VMConfigRateLimiterBandwidth `json:"bandwidth,omitempty"`
	Ops       *VMConfigRateLimiterOps       `json:"ops,omitempty"`
}

type VMConfigRateLimiterOps struct {
	Size         int `json:"size,omitempty"`
	OneTimeBurst int `json:"one_time_burst,omitempty"`
	RefillTime   int `json:"refill_time,omitempty"`
}

type VMConfigRateLimiterBandwidth struct {
	Size         int `json:"size,omitempty"`
	OneTimeBurst int `json:"one_time_burst,omitempty"`
	RefillTime   int `json:"refill_time,omitempty"`
}

type VMConfigPayload struct {
	Firmware  string `json:"firmware,omitempty"`
	Kernel    string `json:"kernel,omitempty"`
	Cmdline   string `json:"cmdline,omitempty"`
	Initramfs string `json:"initramfs,omitempty"`
}

type VMConfigMemory struct {
	Size           string                `json:"size,omitempty"`
	HotPlugSize    int                   `json:"hotplug_size,omitempty"`
	HotPluggedSize int                   `json:"hotplugged_size,omitempty"`
	Mergeable      bool                  `json:"mergeable,omitempty"`
	HotplugMethod  string                `json:"hotplug_method,omitempty"`
	Shared         bool                  `json:"shared,omitempty"`
	HugePages      bool                  `json:"hugepages,omitempty"`
	HugePageSize   int                   `json:"hugepage_size,omitempty"`
	Prefault       bool                  `json:"prefault,omitempty"`
	THP            bool                  `json:"thp,omitempty"`
	Zones          []*VMConfigMemoryZone `json:"zones,omitempty"`
}

type VMConfigMemoryZone struct {
	ID             string `json:"id,omitempty"`
	Size           string `json:"size,omitempty"`
	File           string `json:"file,omitempty"`
	Mergeable      bool   `json:"mergeable,omitempty"`
	Shared         bool   `json:"shared,omitempty"`
	HugePages      bool   `json:"hugepages,omitempty"`
	HugePageSize   int    `json:"hugepage_size,omitempty"`
	HostNumaNode   int    `json:"host_numa_node,omitempty"`
	HotPlugSize    int    `json:"hotplug_size,omitempty"`
	HotPluggedSize int    `json:"hotplugged_size,omitempty"`
	Prefault       bool   `json:"prefault,omitempty"`
}

type VMConfigCPUs struct {
	BootVCPUs    int                     `json:"boot_vcpus,omitempty"`
	MaxVCPUs     int                     `json:"max_vcpus,omitempty"`
	Topology     *VMConfigCPUsTopology   `json:"topology,omitempty"`
	KVMHyperV    bool                    `json:"kvm_hyperv,omitempty"`
	MaxxPhysBits int                     `json:"max_phys_bits,omitempty"`
	Affinity     []*VMConfigCPUsAffinity `json:"affinity,omitempty"`
	Features     map[string]bool         `json:"features,omitempty"`
}

type VMConfigCPUsAffinity struct {
	VCPU     int   `json:"vcpu,omitempty"`
	HostCPUs []int `json:"host_cpus,omitempty"`
}

type VMConfigCPUsTopology struct {
	ThreadsPerCore int `json:"threads_per_core,omitempty"`
	CoresPerDie    int `json:"cores_per_die,omitempty"`
	DiesPerPackage int `json:"dies_per_package,omitempty"`
	Packages       int `json:"packages,omitempty"`
}

type VMDeviceItem struct {
	ID        string            `json:"id,omitempty"`
	Resources []json.RawMessage `json:"resources,omitempty"`
	Children  []string          `json:"children,omitempty"`
	PCIBDF    string            `json:"pci_bdf,omitempty"`
}
