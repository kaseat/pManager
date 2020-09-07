package postgres

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/kaseat/pManager/models"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// GetShares gets shares
func (db Db) GetShares(pid string, onDate string) ([]models.Share, error) {
	i64, err := strconv.ParseInt(pid, 10, 32)
	if err != nil {
		return nil, err
	}
	dtime := time.Now()
	if t, err := time.Parse("2006-01-02T15:04:05Z07:00", onDate); err == nil {
		dtime = t
	}
	bef := dtime.AddDate(0, 0, -5)
	defer timeTrack(time.Now(), "GetShares")
	q := `
select
	x.isin,
	x.vol,
	x.buy_sum,
	x.current_price,
	s.ticker,
	s.title 
from (
	select
		x.isin,
		x.vol,
		x.price as buy_sum,
		coalesce(p.price, 0) as current_price,
		row_number() over (partition by p.isin order by p.ondate desc) as rn
	from (
		select
			x.isin,
			sum(x.vol) as vol,
			sum(x.price) as price
		from (
			select
				o.isin,
				case
					when o.op_id in (2, 10)
						then o.vol * -1
					when o.op_id = 1
						then o.vol
					else 0
					end as vol,
				case
					when o.op_id in (2, 5, 9, 10)
						then o.vol * o.price
					else
						o.vol * o.price * -1
					end as price
			from operations o
			where o.pid = $1
				and o.time <= $2
		) x
		group by x.isin
	) x
	left join prices p
		on p.isin = x.isin
			and p.ondate <= $2
			and p.ondate > $3
	where (x.vol = 0 and x.isin = 'RUB') or x.vol <> 0
	group by x.isin, x.vol, x.price, p.isin, p.ondate
) x
inner join securities s
	on s.isin = x.isin
where rn = 1
	`

	rows, err := db.connection.Query(context.Background(), q, int32(i64), dtime, bef)
	if err != nil {
		return nil, err
	}
	shares := []models.Share{}
	sum := 0.0
	for rows.Next() {
		var isin string
		var vol int64
		var buySum float64
		var currPrice float64
		var ticker string
		var title string

		rows.Scan(&isin, &vol, &buySum, &currPrice, &ticker, &title)

		sum += buySum

		if isin == "RUB         " {
			continue
		}
		share := models.Share{
			ISIN:   isin,
			Ticker: ticker,
			Date:   dtime,
			Volume: vol,
			Price:  currPrice,
		}
		shares = append(shares, share)
	}
	if len(shares) > 0 || sum > 0 {
		share := models.Share{
			ISIN:   "RUB",
			Ticker: "RUB",
			Date:   dtime,
			Volume: 0,
			Price:  float64(int64(sum*100)) / 100,
		}
		shares = append(shares, share)
	}
	return shares, nil
}
