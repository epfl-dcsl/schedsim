#!/usr/bin/python3

import fire
import os
import sys
import subprocess
from multiprocessing import Process
from enum import Enum, auto
import csv
from pprint import pprint

# Only written for single Q
load_levels = [0.01, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 0.95, 0.99]
metrics = ['Count', 'Stolen', 'AVG', 'STDDev',
           '50th', '90th', '95th', '99th', 'Reqs/time_unit']


def run(topo, mu, gen_type, proc_type, num_cores):
    '''
    mu in us
    '''
    service_time_per_core_us = 1 / mu
    rps_capacity_per_core = 1 / service_time_per_core_us * 1000.0 * 1000.0
    total_rps_capacity = rps_capacity_per_core * num_cores
    injected_rps = [load_lvl * total_rps_capacity for load_lvl in load_levels]
    lambdas = [rps / 1000.0 / 1000.0 for rps in injected_rps]
    res_file = "out.txt"
    with open(res_file, 'w') as f:
        for l in lambdas:
            cmd = f"schedsim --topo={topo} --mu={mu} --genType={gen_type} --procType={proc_type} --lambda={l}"
            print(f"Running... {cmd}")
            subprocess.run(cmd, stdout=f, shell=True)


def out_to_csv():
    results = {}
    with open("out.txt", 'r') as f:
        csv_reader = csv.reader(f, delimiter='\t')
        rate = 0
        next_is_res = False
        for row in csv_reader:
            if len(row) >= 3 and "interarrival_rate" in row[2]:
                rate = row[2].split(":")[1]
                results[rate] = {}
            if next_is_res:
                for i, metric in enumerate(metrics):
                    results[rate][metric] = row[i]
            next_is_res = "Count" == row[0]

    pprint(results)

    with open("out.csv", 'w') as f:
        writer = csv.writer(f, delimiter='\t')
        for rate in results:
            # TODO: let user choose
            writer.writerow(
                [results[rate]['50th'], results[rate]['99th']])


if __name__ == "__main__":
    fire.Fire({
        "run": run,
        "csv": out_to_csv
    })
