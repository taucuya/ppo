import pandas as pd
import matplotlib.pyplot as plt
import sys

data = pd.read_csv(sys.argv[1])

percentiles = [50, 70, 75, 90, 99]
percentile_values = {f'p{p}': round(data['duration (ms)'].quantile(p/100), 2) for p in percentiles}
percentile_df = pd.DataFrame([percentile_values])
percentile_df.to_csv('percentiles.csv', index=False)

plt.figure(figsize=(12, 8))
n, bins, patches = plt.hist(data['duration (ms)'], bins=50, alpha=0.7, color='lightblue', edgecolor='black')
colors = ['red', 'orange', 'yellow', 'green', 'purple']
for i, p in enumerate(percentiles):
    value = percentile_values[f'p{p}']
    plt.axvline(value, color=colors[i], linestyle='--', linewidth=2, label=f'p{p}: {value}ms')
plt.xlabel('Latency (ms)')
plt.ylabel('Number of Requests')
plt.title(f'Latency Distribution with Percentiles\nRPS: {data["rps"].iloc[0]}, Total Requests: {len(data)}')
plt.legend()
plt.grid(True, alpha=0.3)
plt.savefig('./img/percentile_histogram.png', dpi=300, bbox_inches='tight')
plt.close()

print(f"p50: {percentile_values['p50']}ms")
print(f"p70: {percentile_values['p70']}ms")
print(f"p75: {percentile_values['p75']}ms")
print(f"p90: {percentile_values['p90']}ms")
print(f"p99: {percentile_values['p99']}ms")