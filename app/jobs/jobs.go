package jobs

import (
	"fmt"
	"github.com/robfig/cron"
	"go-web-demo/app/jobs/ranking"
)

func Start() {
	// Tips1：每条任务都必须支持多进程执行，以方便以后分布式部署该项目
	// Tips2：任务最好能够缠粉独立进程运行，避免与http互相干扰
	c := cron.New()
	var err error
	// 更新七日排行
	err = c.AddFunc("0 0/50 * * * *", ranking.UpdateSevenDaysRankingList)
	if err != nil {
		fmt.Println(err)
	}

	c.Start()
}
