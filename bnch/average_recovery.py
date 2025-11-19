import pandas as pd
import numpy as np
import glob
import os
import sys

def calculate_recovery_stats(out_dir="recovery_out"):
    search_pattern = os.path.join(out_dir, "**", "recovery_time*.csv")
    files = glob.glob(search_pattern, recursive=True)
    
    if not files:
        print(f"No recovery time files found matching: {search_pattern}")
        return
    
    print(f"Found {len(files)} files:")
    for f in files:
        print(f"  {os.path.basename(f)}")
    print()
    
    all_recovery_times = []
    all_rps_values = []
    
    for file_path in files:
        try:
            df = pd.read_csv(file_path)
            if not df.empty and 'recovery_time_ms' in df.columns:
                recovery_data = df[df['rps'].apply(lambda x: str(x).isdigit())]
                if not recovery_data.empty:
                    recovery_times = recovery_data['recovery_time_ms'].tolist()
                    rps_values = recovery_data['rps'].tolist()
                    
                    all_recovery_times.extend(recovery_times)
                    all_rps_values.extend(rps_values)
                    
                    print(f"Processed: {os.path.basename(file_path)}")
                    for rps, time in zip(rps_values, recovery_times):
                        print(f"  RPS: {rps}, Recovery: {time}ms")
        except Exception as e:
            print(f"Error reading {file_path}: {e}")
    
    if not all_recovery_times:
        print("No valid data found")
        return
    
    recovery_array = np.array(all_recovery_times)
    rps_array = np.array(all_rps_values)
    
    print(f"\n=== OVERALL STATISTICS ===")
    print(f"Total samples: {len(recovery_array)}")
    print(f"Average recovery time: {np.mean(recovery_array):.2f}ms")
    print(f"Median recovery time: {np.median(recovery_array):.2f}ms")
    print(f"Minimum recovery time: {np.min(recovery_array):.2f}ms")
    print(f"Maximum recovery time: {np.max(recovery_array):.2f}ms")
    print(f"Standard deviation: {np.std(recovery_array):.2f}ms")
    
    unique_rps = np.unique(rps_array)
    
    print(f"\n=== STATISTICS BY RPS ===")
    results_by_rps = []
    
    for rps in unique_rps:
        mask = rps_array == rps
        rps_times = recovery_array[mask]
        
        stats = {
            'rps': rps,
            'samples': len(rps_times),
            'average_ms': np.mean(rps_times),
            'median_ms': np.median(rps_times),
            'min_ms': np.min(rps_times),
            'max_ms': np.max(rps_times),
            'std_ms': np.std(rps_times),
            'p95_ms': np.percentile(rps_times, 95)
        }
        results_by_rps.append(stats)
        
        print(f"\nRPS {rps}:")
        print(f"  Samples:     {stats['samples']}")
        print(f"  Average:     {stats['average_ms']:.2f}ms")
        print(f"  Median:      {stats['median_ms']:.2f}ms")
        print(f"  Min:         {stats['min_ms']:.2f}ms")
        print(f"  Max:         {stats['max_ms']:.2f}ms")
        print(f"  Std Dev:     {stats['std_ms']:.2f}ms")
        print(f"  95th %ile:   {stats['p95_ms']:.2f}ms")
    
    results_df = pd.DataFrame(results_by_rps)
    output_file = os.path.join(out_dir, "recovery_statistics_detailed.csv")
    results_df.to_csv(output_file, index=False)
    
    summary_file = os.path.join(out_dir, "recovery_summary.txt")
    with open(summary_file, 'w') as f:
        f.write("=== RECOVERY TIME STATISTICS ===\n\n")
        f.write(f"Overall average: {np.mean(recovery_array):.2f}ms\n")
        f.write(f"Overall median: {np.median(recovery_array):.2f}ms\n\n")
        
        for stats in results_by_rps:
            f.write(f"RPS {stats['rps']}:\n")
            f.write(f"  Samples:    {stats['samples']}\n")
            f.write(f"  Average:    {stats['average_ms']:.2f}ms\n")
            f.write(f"  Median:     {stats['median_ms']:.2f}ms\n")
            f.write(f"  Range:      {stats['min_ms']:.2f}ms - {stats['max_ms']:.2f}ms\n")
            f.write(f"  Std Dev:    {stats['std_ms']:.2f}ms\n")
            f.write(f"  95th %ile:  {stats['p95_ms']:.2f}ms\n\n")
    
    print(f"\nDetailed results saved to: {output_file}")
    print(f"Summary saved to: {summary_file}")
    
    return results_by_rps

if __name__ == "__main__":
    out_dir = sys.argv[1] if len(sys.argv) > 1 else "recovery_out"
    calculate_recovery_stats(out_dir)