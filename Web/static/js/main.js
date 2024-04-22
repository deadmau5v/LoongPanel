    function getCPUPercent() {

        axios.get('/api/v1/cpu_percent')
            .then(function (response) {
                // 处理响应数据
                let data = response.data;
                if (data.status === 0) {
                    if (data.percent < 1) {
                        data.percent = 1
                    }
                    CPUCharts1option.series[0].data[0].value = data.percent;
                    CPUCharts1option.series[0].data[1].value = 100 - data.percent;
                    CPUCharts1.setOption(CPUCharts1option);
                }
            })
            .catch(function (error) {
                // 处理错误
                console.log(error);
            });

    }


    function getMemoryPercent() {

        axios.get('/api/v1/ram_percent')
            .then(function (response) {
                // 处理响应数据
                let data = response.data;
                if (data.status === 0) {
                    MemoryCharts1Option.series[0].data[0].value = (data.percent / 1024 / 1024 / 1024);
                    MemoryCharts1Option.series[0].data[1].value = (data.max / 1024 / 1024 / 1024) - (data.percent / 1024 / 1024 / 1024);
                    MemoryCharts1.setOption(MemoryCharts1Option);
                }
            })
            .catch(function (error) {
                // 处理错误
                console.log(error);
            });
    }