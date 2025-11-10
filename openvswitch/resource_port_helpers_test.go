package openvswitch

import (
	"testing"

	"github.com/digitalocean/go-openvswitch/ovs"
)

func TestGetPortAction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected ovs.PortAction
	}{
		{
			name:     "up action",
			input:    "up",
			expected: ovs.PortActionUp,
		},
		{
			name:     "down action",
			input:    "down",
			expected: ovs.PortActionDown,
		},
		{
			name:     "stp action",
			input:    "stp",
			expected: ovs.PortActionSTP,
		},
		{
			name:     "no-stp action",
			input:    "no-stp",
			expected: ovs.PortActionNoSTP,
		},
		{
			name:     "receive action",
			input:    "receive",
			expected: ovs.PortActionReceive,
		},
		{
			name:     "no-receive action",
			input:    "no-receive",
			expected: ovs.PortActionNoReceive,
		},
		{
			name:     "no-receive-stp action",
			input:    "no-receive-stp",
			expected: ovs.PortActionReceiveSTP,
		},
		{
			name:     "forward action",
			input:    "forward",
			expected: ovs.PortActionForward,
		},
		{
			name:     "no-forward action",
			input:    "no-forward",
			expected: ovs.PortActionNoForward,
		},
		{
			name:     "flood action",
			input:    "flood",
			expected: ovs.PortActionFlood,
		},
		{
			name:     "no-flood action",
			input:    "no-flood",
			expected: ovs.PortActionNoFlood,
		},
		{
			name:     "packet-in action",
			input:    "packet-in",
			expected: ovs.PortActionPacketIn,
		},
		{
			name:     "no-packet-in action",
			input:    "no-packet-in",
			expected: ovs.PortActionNoPacketIn,
		},
		{
			name:     "invalid action defaults to up",
			input:    "invalid",
			expected: ovs.PortActionUp,
		},
		{
			name:     "empty action defaults to up",
			input:    "",
			expected: ovs.PortActionUp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPortAction(tt.input)
			if result != tt.expected {
				t.Errorf("GetPortAction(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetPortActionCaseSensitivity(t *testing.T) {
	// Test that the function is case-sensitive (current behavior)
	tests := []struct {
		name  string
		input string
	}{
		{name: "uppercase UP", input: "UP"},
		{name: "mixed case Up", input: "Up"},
		{name: "mixed case Forward", input: "Forward"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetPortAction(tt.input)
			// All case variations should default to Up since switch is case-sensitive
			if result != ovs.PortActionUp {
				t.Errorf("GetPortAction(%q) = %v, want %v (default for unmatched case)",
					tt.input, result, ovs.PortActionUp)
			}
		})
	}
}
