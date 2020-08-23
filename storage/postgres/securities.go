package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/pgtype"
	"github.com/kaseat/pManager/models"
)

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// GetShares gets shares
func (db Db) GetShares(pid string, onDate string) ([]models.Share, error) {
	defer timeTrack(time.Now(), "GetShares")
	q := `
	select
		x.isin,
		x.vol,
		x.price,
		max(p.price) as onpr,
		sum(x.price) over() as pa
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
					when o.op_id in (2, 5, 9 ,10)
						then o.vol * o.price
					else
						o.vol * o.price * -1
					end as price
			from operations o
			where o.pid = 15048870
				and o.time <= '2020-05-27'
		) x
		group by x.isin
	) x
	left join prices p
		on p.isin = x.isin
			and p.ondate <='2020-05-27'
			and p.ondate >'2020-05-20'
	where (x.vol = 0 and x.isin = 'RUB') or x.vol <> 0
	group by x.isin, x.vol, x.price, p.isin
	`

	rows, err := db.connection.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var isin string
		var vol int
		var price float32
		var onpr pgtype.Float4
		var pa float32
		rows.Scan(&isin, &vol, &price, &onpr, &pa)
		fmt.Println(isin, vol, price, onpr.Float, pa)
	}
	return nil, nil
}
