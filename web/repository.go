package web

const repositoryTemplate = `
<html>

<head>
    <title>Stargazer | {{.entry.Repository}}</title>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.9.3/Chart.bundle.min.js"
        integrity="sha256-TQq84xX6vkwR0Qs1qH5ADkP+MvH0W+9E7TdHJsoIQiM=" crossorigin="anonymous"></script>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/Chart.js/2.9.3/Chart.min.css"
        integrity="sha256-aa0xaJgmK/X74WM224KMQeNQC2xYKwlAt08oZqjeF0E=" crossorigin="anonymous" />
    <style>
        body {
            padding-left: 100px;
            padding-right: 100px;
            display: flex;
            flex-direction: column;
            font-family: Arial, Helvetica, sans-serif;
        }

        .title {
            text-align: center;
            font-size: 50px;
            margin-top: 75px;
            margin-bottom: 50px;
        }

        .content {
            text-align: center;
            font-size: 20px;
            margin-top: 75px;
            margin-bottom: 50px;
        }

        .graph {
            margin: 50px;
        }
    </style>
</head>

<body>
	<div class="title">{{.entry.Repository}}</div>
	{{if eq .entry.Status "requested"}}
	<div class="content">Stats are computing, refresh this page in a few minutes!</div>
	{{else}}
    <div class="graph">
        <canvas id="allStars"></canvas>
    </div>
    <div class="graph">
        <canvas id="starPerDay"></canvas>
    </div>
    <script>
        var stats = JSON.parse("{{.stats_json}}");
        var evolutionLabels = [];
        var evolutionData = [];
        for (var i = 0 ; i < stats.evolution.length ; i++) {
            evolutionLabels.push(stats.evolution[i].date);
            evolutionData.push(stats.evolution[i].count);
        }
        var perDaysLabels = [];
        var perDaysData = [];
        for (var i = 0 ; i < stats.per_days.length ; i++) {
            perDaysLabels.push(stats.per_days[i].date);
            perDaysData.push(stats.per_days[i].count);
        }

        var config = {
            type: 'line',
            data: {
                labels: evolutionLabels,
                datasets: [{
                    label: 'Stars evolution',
                    backgroundColor: 'rgba(54, 162, 235, 0.2)',
                    borderColor: 'rgba(54, 162, 235, 1)',
                    data: evolutionData,
                    borderWidth: 1,
                    fill: false,
                    pointRadius: 0
                }]
            },
            options: {
                responsive: true,
                title: {
                    display: false,
                    text: ''
                },
                tooltips: {
                    mode: 'index',
                    intersect: false,
                },
                hover: {
                    mode: 'nearest',
                    intersect: true
                },
                scales: {
                    xAxes: [{
                        type: 'time',
                        time: {
                            unit: 'day'
                        },
                        display: true,
                        scaleLabel: {
                            display: true,
                            labelString: 'Date'
                        }
                    }],
                    yAxes: [{
                        display: true,
                        scaleLabel: {
                            display: true,
                            labelString: 'Count'
                        }
                    }]
                }
            }
        };

        window.onload = function () {
            var ctx = document.getElementById('allStars').getContext('2d');
            window.myLine = new Chart(ctx, config);
        };

        var ctx2 = document.getElementById('starPerDay').getContext('2d');
        var myChart2 = new Chart(ctx2, {
            type: 'line',
            data: {
                labels: perDaysLabels,
                datasets: [{
                    label: 'Stars per days',
                    backgroundColor: 'rgba(54, 162, 235, 0.2)',
                    borderColor: 'rgba(54, 162, 235, 1)',
                    data: perDaysData,
                    borderWidth: 1,
                    lineTension: 0,
                    fill: false
                }]
            },
            options: {
                responsive: true,
                title: {
                    display: false,
                    text: ''
                },
                tooltips: {
                    mode: 'index',
                    intersect: false,
                },
                hover: {
                    mode: 'nearest',
                    intersect: true
                },
                scales: {
                    xAxes: [{
                        type: 'time',
                        time: {
                            unit: 'day'
                        },
                        display: true,
                        scaleLabel: {
                            display: true,
                            labelString: 'Date'
                        }
                    }],
                    yAxes: [{
                        display: true,
                        scaleLabel: {
                            display: true,
                            labelString: 'Count'
                        }
                    }]
                }
            }
        });
	</script>
	{{end}}
</body>

</html>`
