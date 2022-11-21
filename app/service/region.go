package service

import (
	"encoding/json"
	"os"
	"strings"
)

var RegionService = new(regionService)

type regionService struct{}

type JsonRegion struct {
	ParentCode  string `json:"parent_code"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Type        int    `json:"type"`
	ZipCode     string `json:"zip_code"`
	Pinyin      string `json:"pinyin"`
	FirstLetter string `json:"first_letter"`
}

type Region struct {
	Value       string    `json:"value"`
	Label       string    `json:"label"`
	Type        int       `json:"type"`
	ZipCode     string    `json:"zip_code"`
	Pinyin      string    `json:"pinyin"`
	FirstLetter string    `json:"first_letter"`
	Children    []*Region `json:"children"`
}

var RegionList = []*JsonRegion{}
var RegionTree = []*Region{}

func (s *regionService) List() (list []*JsonRegion, err error) {
	if len(RegionList) > 0 {
		list = RegionList
		return
	}

	file, err := os.Open("./data/region.json")
	if err != nil {
		return
	}

	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		return
	}

	err = json.Unmarshal(buffer, &list)
	if err != nil {
		return
	}

	RegionList = list
	return
}

func (s *regionService) Tree() (tree []*Region, err error) {
	if len(RegionTree) > 0 {
		tree = RegionTree
		return
	}

	region_list, err := s.List()
	if err != nil {
		return
	}

	for _, v := range region_list {
		if v.ParentCode == "" {
			tree = append(tree, &Region{
				Value:       v.Code,
				Label:       v.Name,
				Type:        v.Type,
				ZipCode:     v.ZipCode,
				Pinyin:      v.Pinyin,
				FirstLetter: v.FirstLetter,
			})
		}
	}

	for _, v := range tree {
		for _, v2 := range region_list {
			if v2.ParentCode == v.Value {
				v.Children = append(v.Children, &Region{
					Value:       v2.Code,
					Label:       v2.Name,
					Type:        v2.Type,
					ZipCode:     v2.ZipCode,
					Pinyin:      v2.Pinyin,
					FirstLetter: v2.FirstLetter,
				})
			}
		}

		for _, v2 := range v.Children {
			for _, v3 := range region_list {
				if v2.Value == v3.ParentCode {
					v2.Children = append(v2.Children, &Region{
						Value:       v3.Code,
						Label:       v3.Name,
						Type:        v3.Type,
						ZipCode:     v3.ZipCode,
						Pinyin:      v3.Pinyin,
						FirstLetter: v3.FirstLetter,
					})
				}
			}
		}
	}

	RegionTree = tree
	return
}

func (s *regionService) GetNameById(region_id string) (region_full_name string, err error) {
	region_list, err := s.List()
	if err != nil {
		return
	}

	region_bit_arr := strings.Split(region_id, "")
	if len(region_bit_arr) != 6 {
		return
	}

	region_name := []string{}
	for _, v := range region_list {
		if v.Code == region_bit_arr[0]+region_bit_arr[1]+"0000" {
			region_name = append(region_name, v.Name)
		}
		if v.Code == region_bit_arr[0]+region_bit_arr[1]+region_bit_arr[2]+region_bit_arr[3]+"00" {
			region_name = append(region_name, v.Name)
		}
		if v.Code == region_id {
			region_name = append(region_name, v.Name)
		}
	}
	region_full_name = strings.Join(region_name, " / ")
	return
}
