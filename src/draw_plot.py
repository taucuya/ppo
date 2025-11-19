import matplotlib.pyplot as plt
import pandas as pd
import numpy as np
import glob
from datetime import datetime
import os

def read_percentiles_data(file_path):
    with open(file_path, "r") as file:
        lines = file.readlines()

    p50, p75, p90, p95, p99, y = [], [], [], [], [], []

    for line in lines[1:]:
        parts = line.strip().split(',')
        if len(parts) >= 6:
            p50.append(float(parts[0].replace('ms', '')))
            p75.append(float(parts[1].replace('ms', '')))
            p90.append(float(parts[2].replace('ms', '')))
            p95.append(float(parts[3].replace('ms', '')))
            p99.append(float(parts[4].replace('ms', '')))
            y.append(int(parts[5]))
    
    return p50, p75, p90, p95, p99, y

def read_grad_data(file_path):
    with open(file_path, "r") as file:
        lines = file.readlines()

    X, Y = [], []
    for line in lines[1:]:
        if line.strip():
            parts = line.strip().split(',')
            if len(parts) >= 2:
                x = parts[0].replace('ms', '')
                y = parts[1]
                X.append(float(x))
                Y.append(int(y))
    
    return X, Y

def read_metrics_data(file_path):
    df = pd.read_csv(file_path)
    
    df['datetime'] = pd.to_datetime(df['timestamp'], unit='s')
    start_time = df['datetime'].min()
    df['time_seconds'] = (df['datetime'] - start_time).dt.total_seconds()
    
    df['memory_usage_mb'] = df['memory_usage'].str.replace('MiB', '').astype(float)
    df['memory_limit_gb'] = df['memory_limit'].str.replace('GiB', '').astype(float)
    
    return df

def plot_performance_graphs():
    fig1, ax1 = plt.subplots(figsize=(10, 6))
    
    if os.path.exists("../benchmark-results/grad.csv"):
        try:
            x, y = read_grad_data("../benchmark-results/grad.csv")
            if x and y:
                ax1.plot(y, x, 'b-', linewidth=2, marker='o', markersize=4)
                ax1.set_ylabel("Время ответа P95, мс")
                ax1.set_xlabel("Количество пользователей")
                ax1.set_title("P95 vs Пользователи")
                ax1.grid(True, alpha=0.3)
                plt.tight_layout()
                plt.savefig("../benchmark-results/p95_vs_users.png", dpi=300, bbox_inches='tight')
                plt.show()
        except Exception as e:
            print(f"Ошибка при чтении grad.csv: {e}")

    if os.path.exists("../benchmark-results/percentiles.csv"):
        try:
            p50, p75, p90, p95, p99, y = read_percentiles_data("../benchmark-results/percentiles.csv")
            
            fig2, ax2 = plt.subplots(figsize=(10, 6))
            ax2.plot(y, p50, label='P50', marker='o', markersize=3)
            ax2.plot(y, p75, label='P75', marker='s', markersize=3)
            ax2.plot(y, p90, label='P90', marker='^', markersize=3)
            ax2.plot(y, p95, label='P95', marker='d', markersize=3)
            ax2.plot(y, p99, label='P99', marker='*', markersize=3)
            ax2.set_ylabel("Время ответа, мс")
            ax2.set_xlabel("Количество пользователей")
            ax2.set_title("Персентили времени ответа")
            ax2.legend()
            ax2.grid(True, alpha=0.3)
            plt.tight_layout()
            plt.savefig("../benchmark-results/percentiles.png", dpi=300, bbox_inches='tight')
            plt.show()
        except Exception as e:
            print(f"Ошибка при чтении percentiles.csv: {e}")

    if os.path.exists("../benchmark-results/percentiles.csv") and os.path.exists("../benchmark-results/grad.csv"):
        try:
            p50, p75, p90, p95, p99, y = read_percentiles_data("../benchmark-results/percentiles.csv")
            x, y = read_grad_data("../benchmark-results/grad.csv")
            
            percentiles_data = [p50, p75, p90, p95, p99, x]
            titles = ['P50', 'P75', 'P90', 'P95', 'P99', 'P100']
            
            for i in range(6):
                fig, ax = plt.subplots(figsize=(8, 6))
                if i < len(percentiles_data) and percentiles_data[i]:
                    ax.hist(percentiles_data[i], bins=20, color='red', alpha=0.7, edgecolor='black')
                    ax.set_title(titles[i])
                    ax.set_xlabel("Время, мс")
                    ax.set_ylabel("Частота")
                    plt.tight_layout()
                    plt.savefig(f"../benchmark-results/histogram_{titles[i]}.png", dpi=300, bbox_inches='tight')
                    plt.show()
        except Exception as e:
            print(f"Ошибка при построении гистограмм: {e}")

def plot_resource_usage():
    metric_files = glob.glob("../benchmark-results/metrics_*.csv")
    
    if not metric_files:
        print("Файлы с метриками не найдены")
        return
    
    for metric_file in metric_files:
        scenario_name = os.path.basename(metric_file).replace('metrics_', '').replace('.csv', '')
        df = read_metrics_data(metric_file)
        
        fig1, ax1 = plt.subplots(figsize=(10, 6))
        ax1.plot(df['time_seconds'], df['cpu_percent'], 'r-', linewidth=1, alpha=0.8)
        ax1.set_ylabel("CPU, %")
        ax1.set_xlabel("Время, сек")
        ax1.set_title(f"Использование CPU - {scenario_name}")
        ax1.grid(True, alpha=0.3)
        plt.tight_layout()
        plt.savefig(f"../benchmark-results/cpu_usage_{scenario_name}.png", dpi=300, bbox_inches='tight')
        plt.show()
        
        fig2, ax2 = plt.subplots(figsize=(10, 6))
        ax2.plot(df['time_seconds'], df['memory_percent'], 'g-', linewidth=1, alpha=0.8)
        ax2.set_ylabel("Память, %")
        ax2.set_xlabel("Время, сек")
        ax2.set_title(f"Использование памяти - {scenario_name}")
        ax2.grid(True, alpha=0.3)
        plt.tight_layout()
        plt.savefig(f"../benchmark-results/memory_usage_{scenario_name}.png", dpi=300, bbox_inches='tight')
        plt.show()
        
        fig3, ax3 = plt.subplots(figsize=(10, 6))
        ax3.hist(df['cpu_percent'], bins=30, color='red', alpha=0.7, edgecolor='black')
        ax3.set_xlabel("CPU, %")
        ax3.set_ylabel("Частота")
        ax3.set_title(f"Распределение CPU - {scenario_name}")
        plt.tight_layout()
        plt.savefig(f"../benchmark-results/cpu_histogram_{scenario_name}.png", dpi=300, bbox_inches='tight')
        plt.show()
        
        fig4, ax4 = plt.subplots(figsize=(10, 6))
        ax4.hist(df['memory_percent'], bins=30, color='green', alpha=0.7, edgecolor='black')
        ax4.set_xlabel("Память, %")
        ax4.set_ylabel("Частота")
        ax4.set_title(f"Распределение памяти - {scenario_name}")
        plt.tight_layout()
        plt.savefig(f"../benchmark-results/memory_histogram_{scenario_name}.png", dpi=300, bbox_inches='tight')
        plt.show()
        
        print(f"=== Статистика {scenario_name} ===")
        print(f"CPU: среднее = {df['cpu_percent'].mean():.2f}%, макс = {df['cpu_percent'].max():.2f}%")
        print(f"Память: среднее = {df['memory_percent'].mean():.2f}%, макс = {df['memory_percent'].max():.2f}%")
        print(f"Длительность: {df['time_seconds'].max():.1f} сек")
        print(f"Записей: {len(df)}")

def plot_comparison():
    metric_files = glob.glob("../benchmark-results/metrics_*.csv")
    
    if len(metric_files) < 2:
        return
    
    fig1, ax1 = plt.subplots(figsize=(12, 6))
    for metric_file in metric_files:
        scenario_name = os.path.basename(metric_file).replace('metrics_', '').replace('.csv', '')
        df = read_metrics_data(metric_file)
        ax1.plot(df['time_seconds'], df['cpu_percent'], label=scenario_name, linewidth=2, alpha=0.8)
    
    ax1.set_ylabel("CPU, %")
    ax1.set_xlabel("Время, сек")
    ax1.set_title("Сравнение CPU по сценариям")
    ax1.legend()
    ax1.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig("../benchmark-results/cpu_comparison.png", dpi=300, bbox_inches='tight')
    plt.show()
    
    fig2, ax2 = plt.subplots(figsize=(12, 6))
    for metric_file in metric_files:
        scenario_name = os.path.basename(metric_file).replace('metrics_', '').replace('.csv', '')
        df = read_metrics_data(metric_file)
        ax2.plot(df['time_seconds'], df['memory_percent'], label=scenario_name, linewidth=2, alpha=0.8)
    
    ax2.set_ylabel("Память, %")
    ax2.set_xlabel("Время, сек")
    ax2.set_title("Сравнение памяти по сценариям")
    ax2.legend()
    ax2.grid(True, alpha=0.3)
    plt.tight_layout()
    plt.savefig("../benchmark-results/memory_comparison.png", dpi=300, bbox_inches='tight')
    plt.show()

if __name__ == "__main__":
    print("Анализ результатов нагрузочного тестирования")
    
    os.makedirs("../benchmark-results", exist_ok=True)
    
    print("\n1. Построение графиков производительности...")
    plot_performance_graphs()
    
    print("\n2. Анализ использования CPU и RAM...")
    plot_resource_usage()
    
    print("\n3. Сравнительный анализ сценариев...")
    plot_comparison()
    
    print("\nАнализ завершен! Графики сохранены в ../benchmark-results/")