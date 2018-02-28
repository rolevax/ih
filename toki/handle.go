package main

import (
	"log"

	"github.com/rolevax/ih/ako/ss"
	"github.com/rolevax/ih/saki"
)

type resp func(interface{})

var tables = map[int64]saki.TableSession{}

func handlePing(msg *ss.TablePing, resp resp) {
	resp(&ss.TableOutputs{})
}

func handleStart(msg *ss.TableStart, resp resp) {
	if _, ok := tables[msg.Tid]; ok {
		log.Fatalln("table already exist")
	}
	if len(msg.Gids) != 4 {
		log.Fatalln("girl length not 4")
	}

	table := saki.NewTableSession(
		int(msg.Gids[0]), int(msg.Gids[1]),
		int(msg.Gids[2]), int(msg.Gids[3]),
	)
	tables[msg.Tid] = table

	mails := table.Start()
	defer saki.DeleteMailVector(mails)

	log.Println("new", msg.Tid)
	output(msg.Tid, mails, resp)
}

func handleAction(msg *ss.TableAction, resp resp) {
	table, ok := tables[msg.Tid]
	if !ok {
		outputTableTan90(msg.Tid, resp)
		return
	}

	mails := table.Action(
		int(msg.Who),
		msg.ActStr,
		int(msg.ActArg),
		msg.ActTile,
		int(msg.Nonce),
	)
	defer saki.DeleteMailVector(mails)

	output(msg.Tid, mails, resp)
}

func handleSweepOne(msg *ss.TableSweepOne, resp resp) {
	table, ok := tables[msg.Tid]
	if !ok {
		outputTableTan90(msg.Tid, resp)
		return
	}

	mails := table.SweepOne(int(msg.Who))
	defer saki.DeleteMailVector(mails)

	output(msg.Tid, mails, resp)
}

func handleSweepAll(msg *ss.TableSweepAll, resp resp) {
	table, ok := tables[msg.Tid]
	if !ok {
		outputTableTan90(msg.Tid, resp)
		return
	}

	mails := table.SweepAll()
	defer saki.DeleteMailVector(mails)

	output(msg.Tid, mails, resp)
}

func handleDeleteIfAny(msg *ss.TableDeleteIfAny, resp resp) {
	deleteTableIfAny(msg.Tid)

	resp(&ss.TableOutputs{
		Tid:   msg.Tid,
		Mails: nil,
	})
}

func output(tid int64, mv saki.MailVector, resp resp) {
	reply := &ss.TableOutputs{
		Tid:   tid,
		Mails: makeMails(mv),
	}

	resp(reply)
}

func outputTableTan90(tid int64, resp resp) {
	reply := &ss.TableOutputs{
		Tid: tid,
		Mails: []*ss.TableMail{
			&ss.TableMail{
				Who:     -1,
				Content: `{"Type": "table-tan90"}`,
			},
		},
	}

	resp(reply)
}

func makeMails(mails saki.MailVector) []*ss.TableMail {
	res := []*ss.TableMail{}

	size := int(mails.Size())
	for i := 0; i < size; i++ {
		res = append(res, &ss.TableMail{
			Who:     int64(mails.Get(i).GetTo()),
			Content: mails.Get(i).GetMsg(),
		})
	}

	return res
}

func deleteTableIfAny(tid int64) {
	table, ok := tables[tid]
	if ok {
		saki.DeleteTableSession(table)
		delete(tables, tid)
		log.Println("delete", tid)
	}
}
