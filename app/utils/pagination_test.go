package utils

import "testing"

func TestPaginateByPageNumSize(t *testing.T) {
	p := PaginateByPageNumSize(100, 0, 10)
	if p.PagesCount != 10 || p.HasNext != true || p.NextPageNum != 1 || p.HasPrev != false || p.PrevPageNum != -1 {
		t.Fatalf("paginate error %+v", p)
	}
	p = PaginateByPageNumSize(100, 1, 10)
	if p.PagesCount != 10 || p.HasNext != true || p.NextPageNum != 2 || p.HasPrev != false || p.PrevPageNum != 0 {
		t.Fatalf("paginate error %+v", p)
	}
	p = PaginateByPageNumSize(100, 3, 10)
	if p.PagesCount != 10 || p.HasNext != true || p.NextPageNum != 4 || p.HasPrev != true || p.PrevPageNum != 2 {
		t.Fatalf("paginate error %+v", p)
	}
	p = PaginateByPageNumSize(100, 10, 10)
	if p.PagesCount != 10 || p.HasNext != false || p.NextPageNum != 11 || p.HasPrev != true || p.PrevPageNum != 9 {
		t.Fatalf("paginate error %+v", p)
	}
	p = PaginateByPageNumSize(100, 11, 10)
	if p.PagesCount != 10 || p.HasNext != false || p.NextPageNum != 12 || p.HasPrev != true || p.PrevPageNum != 10 {
		t.Fatalf("paginate error %+v", p)
	}
	p = PaginateByPageNumSize(100, 3, 13)
	if p.PagesCount != 8 || p.HasNext != true || p.NextPageNum != 4 || p.HasPrev != true || p.PrevPageNum != 2 {
		t.Fatalf("paginate error %+v", p)
	}
	p = PaginateByPageNumSize(100, -1, -1)
	if p.PagesCount != 0 || p.HasNext != false || p.NextPageNum != 0 || p.HasPrev != false || p.PrevPageNum != -2 {
		t.Fatalf("paginate error %+v", p)
	}
	p = PaginateByPageNumSize(0, -1, -1)
	if p != nil {
		t.Fatalf("paginate should nil")
	}
}

func TestPaginateByOffsetLimit(t *testing.T) {
	p := PaginateByOffsetLimit(100, 0, 10)
	if p.HasNext != true || p.HasPrev != false || p.NextPageNum != 2 || p.PagesCount != 10 || p.PageNum != 1 || p.PageSize != 10 || p.PrevPageNum != 0 || p.TotalCount != 100 {
		t.Fatalf("paginate error %+v", p)
	}
	p = PaginateByOffsetLimit(100, 10, 10)
	if p.HasNext != true || p.HasPrev != true || p.NextPageNum != 3 || p.PagesCount != 10 || p.PageNum != 2 || p.PageSize != 10 || p.PrevPageNum != 1 || p.TotalCount != 100 {
		t.Fatalf("paginate error %+v", p)
	}
	p = PaginateByOffsetLimit(100, 90, 10)
	if p.HasNext != false || p.HasPrev != true || p.NextPageNum != 11 || p.PagesCount != 10 || p.PageNum != 10 || p.PageSize != 10 || p.PrevPageNum != 9 || p.TotalCount != 100 {
		t.Fatalf("paginate error %+v", p)
	}
	p = PaginateByOffsetLimit(100, 100, 10)
	if p.HasNext != false || p.HasPrev != true || p.NextPageNum != 12 || p.PagesCount != 10 || p.PageNum != 11 || p.PageSize != 10 || p.PrevPageNum != 10 || p.TotalCount != 100 {
		t.Fatalf("paginate error %+v", p)
	}
	p = PaginateByOffsetLimit(100, 110, 10)
	if p.HasNext != false || p.HasPrev != true || p.NextPageNum != 13 || p.PagesCount != 10 || p.PageNum != 12 || p.PageSize != 10 || p.PrevPageNum != 11 || p.TotalCount != 100 {
		t.Fatalf("paginate error %+v", p)
	}
}
