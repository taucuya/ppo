import pandas as pd
import matplotlib.pyplot as plt
import sys
import os
from datetime import datetime

def convert_memory_to_mb(memory_str):
    """Convert memory string to MB for plotting"""
    if pd.isna(memory_str) or memory_str == '0':
        return 0
    
    memory_str = str(memory_str).strip()
    try:
        if 'GiB' in memory_str:
            return float(memory_str.replace('GiB', '').strip()) * 1024
        elif 'MiB' in memory_str:
            return float(memory_str.replace('MiB', '').strip())
        elif 'KiB' in memory_str:
            return float(memory_str.replace('KiB', '').strip()) / 1024
        else:
            return float(memory_str)
    except (ValueError, TypeError):
        return 0

def plot_http_metrics(http_csv_file, output_plot):
    http_df = pd.read_csv(http_csv_file)
    print(f"HTTP metrics: {len(http_df)} data points from {http_csv_file}")
    
    if len(http_df) == 0:
        print("No data found in HTTP CSV file")
        return
    
    def convert_timestamp(df):
        if 'timestamp' in df.columns:
            df['datetime'] = pd.to_datetime(df['timestamp'], unit='ms')
            df['seconds'] = (df['timestamp'] - df['timestamp'].min()) / 1000
        return df
    
    http_df = convert_timestamp(http_df)

    # Convert memory_usage to MB for plotting
    if 'memory_usage' in http_df.columns:
        http_df['memory_usage_mb'] = http_df['memory_usage'].apply(convert_memory_to_mb)
        print(f"Memory usage range: {http_df['memory_usage_mb'].min():.1f} - {http_df['memory_usage_mb'].max():.1f} MB")

    fig, (ax1, ax2) = plt.subplots(1, 2, figsize=(18, 6))
    fig.suptitle('Metrics: HTTP', fontsize=16, fontweight='bold')
    
    http_cpu_label = 'HTTP'
    http_mem_label = 'HTTP'
    
    # Plot CPU usage
    if 'cpu_percent' in http_df.columns:
        http_cpu_avg = http_df['cpu_percent'].mean()
        http_cpu_min = http_df['cpu_percent'].min()
        http_cpu_max = http_df['cpu_percent'].max()
        http_cpu_label = f'HTTP (avg={http_cpu_avg:.1f}%, min={http_cpu_min:.1f}%, max={http_cpu_max:.1f}%)'
        
        ax1.plot(http_df['seconds'], http_df['cpu_percent'], 
                'b-', linewidth=2, label=http_cpu_label, alpha=0.8)
    
    ax1.set_xlabel('Time (seconds)')
    ax1.set_ylabel('CPU Usage (%)')
    ax1.set_title('CPU Usage Over Time')
    ax1.grid(True, alpha=0.3)
    ax1.legend(fontsize=9)
    
    # Plot Memory usage in MB
    if 'memory_usage_mb' in http_df.columns:
        http_mem_avg = http_df['memory_usage_mb'].mean()
        http_mem_min = http_df['memory_usage_mb'].min()
        http_mem_max = http_df['memory_usage_mb'].max()
        http_mem_label = f'HTTP (avg={http_mem_avg:.1f}MB, min={http_mem_min:.1f}MB, max={http_mem_max:.1f}MB)'
        
        ax2.plot(http_df['seconds'], http_df['memory_usage_mb'], 
                'r-', linewidth=2, label=http_mem_label, alpha=0.8)
        
        ax2.set_ylabel('Memory Usage (MB)')
    elif 'memory_usage' in http_df.columns:
        # Fallback: plot original memory_usage values if conversion failed
        http_mem_avg = http_df['memory_usage'].mean()
        http_mem_min = http_df['memory_usage'].min()
        http_mem_max = http_df['memory_usage'].max()
        http_mem_label = f'HTTP (avg={http_mem_avg:.1f}, min={http_mem_min:.1f}, max={http_mem_max:.1f})'
        
        ax2.plot(http_df['seconds'], http_df['memory_usage'], 
                'r-', linewidth=2, label=http_mem_label, alpha=0.8)
        
        ax2.set_ylabel('Memory Usage')
    
    ax2.set_xlabel('Time (seconds)')
    ax2.set_title('Memory Usage Over Time')
    ax2.grid(True, alpha=0.3)
    ax2.legend(fontsize=9)
    
    plt.tight_layout()
    
    plt.savefig(output_plot, dpi=300, bbox_inches='tight')
    print(f"HTTP metrics plot saved to: {output_plot}")
    
    print("\nHTTP Metrics:")
    if 'cpu_percent' in http_df.columns:
        cpu_stats = http_df['cpu_percent'].describe()
        print(f"  CPU Usage: avg={cpu_stats['mean']:.1f}%, min={cpu_stats['min']:.1f}%, max={cpu_stats['max']:.1f}%")
    
    if 'memory_usage_mb' in http_df.columns:
        mem_stats = http_df['memory_usage_mb'].describe()
        print(f"  Memory Usage: avg={mem_stats['mean']:.1f}MB, min={mem_stats['min']:.1f}MB, max={mem_stats['max']:.1f}MB")
        print(f"  Memory Usage (original): {http_df['memory_usage'].iloc[0]} - {http_df['memory_usage'].iloc[-1]}")
    elif 'memory_usage' in http_df.columns:
        mem_stats = http_df['memory_usage'].describe()
        print(f"  Memory Usage: avg={mem_stats['mean']:.1f}, min={mem_stats['min']:.1f}, max={mem_stats['max']:.1f}")

def main():    
    http_csv_file = sys.argv[1]
    output_plot = sys.argv[2]
    
    if not os.path.exists(http_csv_file):
        print(f"HTTP metrics CSV file not found: {http_csv_file}")
        sys.exit(1)
    
    os.makedirs(os.path.dirname(output_plot) if os.path.dirname(output_plot) else '.', exist_ok=True)
    
    plot_http_metrics(http_csv_file, output_plot)

if __name__ == "__main__":
    main()