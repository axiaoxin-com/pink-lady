package utils

import "testing"

func TestPaginate(t *testing.T) {
	p := Paginate(100, 1, 10)
	if p.PagesCount != 10 || p.HasNext != true || p.NextPageNum != 2 || p.HasPrev != false || p.PrevPageNum != 1 {
		t.Errorf("paginate error %v", p)
	}
	p = Paginate(100, 3, 10)
	if p.PagesCount != 10 || p.HasNext != true || p.NextPageNum != 4 || p.HasPrev != true || p.PrevPageNum != 2 {
		t.Errorf("paginate error %v", p)
	}
	p = Paginate(100, 10, 10)
	if p.PagesCount != 10 || p.HasNext != false || p.NextPageNum != 10 || p.HasPrev != true || p.PrevPageNum != 9 {
		t.Errorf("paginate error %v", p)
	}
	p = Paginate(100, 3, 13)
	if p.PagesCount != 8 || p.HasNext != true || p.NextPageNum != 4 || p.HasPrev != true || p.PrevPageNum != 2 {
		t.Errorf("paginate error %v", p)
	}
	p = Paginate(100, -1, -1)
	if p.PagesCount != 1 || p.HasNext != false || p.NextPageNum != 1 || p.HasPrev != false || p.PrevPageNum != 1 {
		t.Errorf("paginate error %v", p)
	}
	p = Paginate(0, -1, -1)
	if p.PagesCount != 0 || p.HasNext != false || p.NextPageNum != 0 || p.HasPrev != false || p.PrevPageNum != 0 {
		t.Errorf("paginate error %v", p)
	}
}
