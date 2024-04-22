axios.get('/api/v1/disks')
    .then(function (response) {
        // 处理响应数据
        let data = response.data;
        let ul = document.getElementById('disk_list');
        for (const disk of data.disks) {
            let li = document.createElement("li")
            let span1 = document.createElement("span")
            let span2 = document.createElement("span")
            let progress = document.createElement("div")
            span1.innerText = disk.FileSystem
            span2.innerText = disk.MountedPath
            progress.className = "progress"
            let fill = document.createElement("div")
            fill.className = "fill"
            let used = disk.UsedMemory / disk.MaxMemory * 100
            fill.style.width = used + "%"

            let span3 = document.createElement("span")
            let max = disk.MaxMemory / 1024 / 1024 / 1024
            let use = disk.UsedMemory / 1024 / 1024 / 1024
            max = max.toFixed(2)
            use = use.toFixed(2)
            if (max > 1024) {
                max = (max / 1024).toFixed(2) + "TB"
            } else {
                max = max + "GB"
            }
            if (use > 1024) {
                use = (use / 1024).toFixed(2) + "TB"
            } else {
                use = use + "GB"
            }
            span3.innerText = use + "/" + max
            span3.style.width = "30%"
            span3.style.textAlign = "left"
            span3.style.marginLeft = "20px"
            progress.appendChild(fill)
            li.appendChild(span1)
            li.appendChild(span2)
            li.appendChild(progress)
            li.appendChild(span3)
            ul.appendChild(li)
            // Todo 把 use 和 max 分开到两个span 更好看

        }
    })
    .catch(function (error) {
        // 处理错误
        console.log(error);
    });


let CPUCharts1 = echarts.init(document.getElementById('CPUCharts1'));
let CPUCharts1option;
CPUCharts1option = {
    animation: true,
    legend: {
        show: false
    },
    series: [
        {
            name: 'CPU占用率',
            type: 'pie',
            radius: ['45%', '50%'],
            itemStyle: {
                borderRadius: 10,
                borderWidth: 5
            },
            data: [
                {value: 1, name: '占用'},
                {value: 99, name: '空闲'}
            ],
            color: [
                '#da4343',
                '#bf9595'
            ],
            startAngle: 270,
        }

    ],
    label: {
        show: true,
        position: 'center',
        fontSize: '20',
        fontWeight: 'bold',
        color: 'black'
    }
};
CPUCharts1.setOption(CPUCharts1option);
CPUCharts1option.animation = false
setInterval(getCPUPercent, 1100);

let MemoryCharts1 = echarts.init(document.getElementById('MemoryCharts1'))
let MemoryCharts1Option
MemoryCharts1Option = {
    animation: true,
    legend: {
        show: false
    },
    series: [
        {
            name: '内存占用率',
            type: 'pie',
            radius: ['45%', '50%'],
            itemStyle: {
                borderRadius: 10,
                borderWidth: 5
            },
            data: [
                {value: 1, name: '占用'},
                {value: 99, name: '空闲'}
            ],
            color: [
                '#da4343',
                '#bf9595'
            ],
            startAngle: 270,
        }

    ],
    label: {
        show: true,
        position: 'center',
        fontSize: '20',
        fontWeight: 'bold',
        color: 'black'
    }
};
MemoryCharts1.setOption(MemoryCharts1Option);
MemoryCharts1Option.animation = false
getMemoryPercent()
setInterval(getMemoryPercent, 10500); // 避免并发