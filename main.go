package main

import (
	"net/http"
	"sort"
	"strconv"
	"github.com/gorilla/mux"
	"strings"
	"io/ioutil"
	"log"
	"encoding/json"
)

type Pet struct {
	ID        int    `json:"id"`
	Type      string `json:"type"`
	Breed     string `json:"breed"`
	Name      string `json:"name"`
	Character string `json:"character"`
	Favorite  string `json:"favorite"`
}

type Hasil struct {
	Name   string
	Length int
}

var pets = []Pet{
	{1, "Anjing", "Golden Retriever", "Otto", "Energik dan senang bermain bola", "ya"},
	{2, "Anjing", "Siberian Husky", "Max", "Bulu lebat dan mata biru", "ya"},
	{3, "Anjing", "Beagle", "Bob", "Ceria dan aktif mengajak esa bermain di taman", "tidak"},
	{4, "Kucing", "Persia", "Luna", "Anggun dan manja", "ya"},
	{5, "Kucing", "British Short Hair", "Milo", "Cerdas dan aktif", "ya"},
	{6, "Ikan", "Koi", "Nana", "Ikan Koi yang indah", "tidak"},
	{7, "Ikan", "Mas", "Goldie", "Ikan mas berwarna cerah", "tidak"},
}

// Mengecek apakah string adalah palindrome (tanpa membedakan huruf besar dan kecil)
func isPalindrome(s string) bool {
	s = strings.ToLower(s)
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		if runes[i] != runes[n-1-i] {
			return false
		}
	}
	return true
}

// Fungsi utama untuk mendapatkan hewan peliharaan dengan nama palindrome
func CekHewanPalindrome(pets []Pet) []Hasil {
	var hasil []Hasil
	for _, pet := range pets {
		if isPalindrome(pet.Name) {
			hasil = append(hasil, Hasil{
				Name:   pet.Name,
				Length: len([]rune(pet.Name)), // untuk menghitung panjang Unicode dengan benar
			})
		}
	}
	return hasil
}

func palindromeHandler(w http.ResponseWriter, r *http.Request){

    result := CekHewanPalindrome(pets)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func getPets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pets)
}

func addPet(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Hanya method POST yang diizinkan", http.StatusMethodNotAllowed)
		return
	}

    last_id := 0
    if len(pets) > 0 {
        last_id = pets[len(pets)-1].ID
    }

    var new_pet Pet
    err := json.NewDecoder(r.Body).Decode(&new_pet)
    if err != nil {
    		http.Error(w, "Invalid JSON", http.StatusBadRequest)
    		return
    }
    defer r.Body.Close()
    new_pet.ID = last_id + 1
    pets = append(pets, new_pet)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(new_pet)
}

func updatePet(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
        http.Error(w, "ID tidak valid", http.StatusBadRequest)
        return
	}

	var updatedFields map[string]string
    if err := json.NewDecoder(r.Body).Decode(&updatedFields); err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

    for i, pet := range pets {
        if pet.ID == id {
            // Update hanya field yang ada di JSON body
            if name, ok := updatedFields["name"]; ok {
                pets[i].Name = name
            }
            if breed, ok := updatedFields["breed"]; ok {
                        pets[i].Breed = breed
            }
            if character, ok := updatedFields["character"]; ok {
                pets[i].Character = character
            }
            if favorite, ok := updatedFields["favorite"]; ok {
                pets[i].Favorite = favorite
            }
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(pets[i])
            return
        }
    }
    http.Error(w, "Pet tidak ditemukan", http.StatusNotFound)
}


func sortByNameAscending(data []Pet) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].Name < data[j].Name
	})
}

func sortByNameDescending(data []Pet) {
	sort.Slice(data, func(i, j int) bool {
		return data[i].Name > data[j].Name
	})
}

func sortHandler(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter ?order=
	order := r.URL.Query().Get("order")
    // 	tampilkan hanya untuk favorite
    var FavoritePets []Pet
	for _, pet := range pets {
        if pet.Favorite == "ya" {
            FavoritePets = append(FavoritePets, pet)
        }
    }
	if order == "asc" {
        sortByNameAscending(FavoritePets)
	}else{
        sortByNameDescending(FavoritePets)
	}
     w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(FavoritePets)
}

func petHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addPet(w, r)
	case http.MethodGet:
		getPets(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func petCount (w http.ResponseWriter, r *http.Request){
       // Ambil query parameter "typepet"
          typepet := r.URL.Query().Get("typepet")

          // Tampilkan hanya untuk type yang cocok
          var TypePets []Pet
          for _, pet := range pets {
              if strings.ToLower(pet.Type) == strings.ToLower(typepet) {
                  TypePets = append(TypePets, pet)
              }
          }

          // Membuat response dalam bentuk JSON
          result := map[string]interface{}{
              "total": len(TypePets),
              "data":  TypePets,
          }

          // Set header Content-Type menjadi application/json
          w.Header().Set("Content-Type", "application/json")

          // Encode result ke dalam response body
          json.NewEncoder(w).Encode(result)
}

func sumArray (w http.ResponseWriter, r *http.Request){
       // Ambil query parameter "typepet"
            data := []int{15,18,3,9,6,2,12,14};
          // Tampilkan hanya untuk type yang cocok
          var genap []int
          var total = 0
          for _, item := range data {
              if  item%2 == 0 {
                  total = total + item
                  genap = append(genap, item)
              }
          }

          result := map[string]interface{}{
              "total": total,
              "data":  genap,
          }

          // Set header Content-Type menjadi application/json
          w.Header().Set("Content-Type", "application/json")

          // Encode result ke dalam response body
          json.NewEncoder(w).Encode(result)
}

func isAnagramHandler(w http.ResponseWriter, r *http.Request) {
    	str1 := []string{"kamu", "buta", "dia"}
    	str2 := []string{"muka", "buat", "aku"}

    	var results []map[string]interface{}

    	for i := 0; i < len(str1); i++ {
    		// Normalisasi string: hilangkan spasi dan ubah ke huruf kecil
    		s1 := strings.ToLower(strings.ReplaceAll(str1[i], " ", ""))
    		s2 := strings.ToLower(strings.ReplaceAll(str2[i], " ", ""))

    		isAnagram := true

    		// Jika panjang tidak sama, langsung bukan anagram
    		if len(s1) != len(s2) {
    			isAnagram = false
    		} else {
    			count := make(map[rune]int)
    			for _, char := range s1 {
    				count[char]++
    			}
    			for _, char := range s2 {
    				count[char]--
    				if count[char] < 0 {
    					isAnagram = false
    					break
    				}
    			}
    		}

    		// Tambahkan hasil ke response
    		results = append(results, map[string]interface{}{
    			"str1":      str1[i],
    			"str2":      str2[i],
    			"isAnagram": isAnagram,
    		})
    	}

    	// Set header dan encode ke JSON
    	w.Header().Set("Content-Type", "application/json")
    	json.NewEncoder(w).Encode(map[string]interface{}{
    		"results": results,
    	})
}

func openFile(w http.ResponseWriter, r *http.Request){
        filename := "case.json"
// Baca isi file JSON
    content, err := ioutil.ReadFile(filename)
    if err != nil {
        log.Fatalf("Gagal membaca file: %v", err)
    }

    // Decode JSON input
    var input Input
    if err := json.Unmarshal(content, &input); err != nil {
        log.Fatalf("Gagal decode JSON: %v", err)
    }

    result := make(map[string]map[string][]ItemBrief)
    categoryTotals := make(map[string]int)
    codeTotals := make(map[string]map[string]int)
    totalAll := 0

    // Proses data
    for _, item := range input.Data {
        totalAll += item.Total

        // Inisialisasi map jika belum ada
        if result[item.Category] == nil {
            result[item.Category] = make(map[string][]ItemBrief)
            codeTotals[item.Category] = make(map[string]int)
        }

        // Simpan data ringkas
        result[item.Category][item.Code] = append(result[item.Category][item.Code], ItemBrief{
            Name:  item.Name,
            Total: item.Total,
        })

        // Hitung total per code dan per category
        codeTotals[item.Category][item.Code] += item.Total
        categoryTotals[item.Category] += item.Total
    }

    // Susun hasil akhir
    var final FinalOutput
    final.Total = totalAll

    // Untuk konsistensi urutan, bisa pakai sort
    var categories []string
    for cat := range result {
        categories = append(categories, cat)
    }
    sort.Strings(categories)

    for _, cat := range categories {
        catOut := CategoryOut{
            Category: cat,
            Total:    categoryTotals[cat],
            Data:     make(map[string]CodeGroupedOutput),
        }

        for code, items := range result[cat] {
            catOut.Data[code] = CodeGroupedOutput{
                Total: codeTotals[cat][code],
                Data:  items,
            }
        }

        final.Data = append(final.Data, catOut)
    }

    // Encode ke JSON hasil akhir
    outputJSON, err := json.MarshalIndent(final, "", "  ")
    if err != nil {
        log.Fatalf("Gagal encode hasil akhir: %v", err)
    }


//         ubah bentuk json
        var parsed interface{}
        json.Unmarshal(outputJSON, &parsed)


        // Set header dan encode ke JSON
    	w.Header().Set("Content-Type", "application/json")
    	json.NewEncoder(w).Encode(map[string]interface{}{
    		"results": parsed,
    	})
}

type Item struct {
    Category string `json:"category"`
    Code     string `json:"code"`
    Name     string `json:"name"`
    Total    int    `json:"total"`
}

type Input struct {
    Data []Item `json:"data"`
}

type FinalOutput struct {
    Total int           `json:"total"`
    Data  []CategoryOut `json:"data"`
}

type CategoryOut struct {
    Category string                       `json:"category"`
    Total    int                          `json:"total"`
    Data     map[string]CodeGroupedOutput `json:"data"`
}

type CodeGroupedOutput struct {
    Total int         `json:"total"`
    Data  []ItemBrief `json:"data"`
}

type ItemBrief struct {
    Name  string `json:"name"`
    Total int    `json:"total"`
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/api/pet/{id}", updatePet).Methods("PATCH")
    r.HandleFunc("/api/pet", petHandler).Methods("GET", "POST")
    r.HandleFunc("/api/pet-by-type", petCount).Methods("GET")
    r.HandleFunc("/api/pet/favorite", sortHandler).Methods("GET")
    r.HandleFunc("/api/pet/sum-array", sumArray).Methods("GET")
    r.HandleFunc("/api/pet/palindrome", palindromeHandler).Methods("GET")
    r.HandleFunc("/api/pet/anagram", isAnagramHandler).Methods("GET")
    r.HandleFunc("/api/pet/reformat-json", openFile).Methods("GET")
	http.ListenAndServe(":8080", r)
}
