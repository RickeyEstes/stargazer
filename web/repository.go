package web

const repositoryTemplate = `
<html>

<head>
    <title>Stargazer | {{.Repository}}</title>
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
	<div class="title">{{.Repository}}</div>
	{{if eq .Status "requested"}}
	<div class="content">Stats are computing, refresh this page in a few minutes!</div>
	{{else}}
    <div class="graph">
        <canvas id="allStars"></canvas>
    </div>
    <div class="graph">
        <canvas id="starPerDay"></canvas>
    </div>
    <script>
        var config = {
            type: 'line',
            data: {
                labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
                datasets: [{
                    label: 'Stars evolution',
                    backgroundColor: 'rgba(54, 162, 235, 0.2)',
                    borderColor: 'rgba(54, 162, 235, 1)',
                    data: [
                        1, 2, 3, 4, 5, 7, 9
                    ],
                    fill: false,
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
            type: 'bar',
            data: {
                labels: ['Red', 'Blue', 'Yellow', 'Green', 'Purple', 'Orange'],
                datasets: [{
                    label: '# of Votes',
                    data: [12, 19, 3, 5, 2, 3],
                    backgroundColor: [
                        'rgba(255, 99, 132, 0.2)',
                        'rgba(54, 162, 235, 0.2)',
                        'rgba(255, 206, 86, 0.2)',
                        'rgba(75, 192, 192, 0.2)',
                        'rgba(153, 102, 255, 0.2)',
                        'rgba(255, 159, 64, 0.2)'
                    ],
                    borderColor: [
                        'rgba(255, 99, 132, 1)',
                        'rgba(54, 162, 235, 1)',
                        'rgba(255, 206, 86, 1)',
                        'rgba(75, 192, 192, 1)',
                        'rgba(153, 102, 255, 1)',
                        'rgba(255, 159, 64, 1)'
                    ],
                    borderWidth: 1
                }]
            },
            options: {
                scales: {
                    yAxes: [{
                        ticks: {
                            beginAtZero: true
                        }
                    }]
                }
            }
        });
	</script>
	{{end}}
</body>

</html>`
