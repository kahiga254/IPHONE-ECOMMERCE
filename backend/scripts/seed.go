package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

const placeholder = "https://placehold.co/600x600/png"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Println("✅ Connected to database")

	// Seed in order
	seedCategories()
	seedApplePhones()
	seedAppleAccessories()
	seedSamsungPhones()
	seedSamsungAccessories()
	seedAdminUser()

	log.Println("✅ Seeding completed successfully")
}

// ─── Helper Functions ─────────────────────────────────────────────────────────

func createCategory(name, slug, description string, parentID *string) string {
	var id string
	err := db.QueryRow(`
		INSERT INTO categories (name, slug, description, image_url)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (slug) DO UPDATE SET name = EXCLUDED.name
		RETURNING id`,
		name, slug, description, placeholder,
	).Scan(&id)
	if err != nil {
		log.Fatalf("Failed to create category %s: %v", name, err)
	}
	log.Printf("   ✓ Category: %s", name)
	return id
}

func createProduct(name, slug, description, categoryID string, basePrice float64, metaTitle string) string {
	var id string
	err := db.QueryRow(`
		INSERT INTO products (name, slug, description, category_id, base_price, meta_title, meta_description, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE)
		ON CONFLICT (slug) DO UPDATE SET base_price = EXCLUDED.base_price
		RETURNING id`,
		name, slug, description, categoryID, basePrice, metaTitle,
		fmt.Sprintf("Buy %s in Kenya. Best prices, genuine product.", name),
	).Scan(&id)
	if err != nil {
		log.Fatalf("Failed to create product %s: %v", name, err)
	}
	return id
}

func createVariant(productID, sku, color, storage string, price float64, stock int) {
	images, _ := json.Marshal([]string{placeholder})
	_, err := db.Exec(`
		INSERT INTO product_variants (product_id, sku, color, storage, price, stock, images)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (sku) DO UPDATE SET price = EXCLUDED.price`,
		productID, sku, color, storage, price, stock, images,
	)
	if err != nil {
		log.Fatalf("Failed to create variant %s: %v", sku, err)
	}
}

func createSpec(productID, key, value string) {
	_, err := db.Exec(`
		INSERT INTO product_specs (product_id, spec_key, spec_value)
		VALUES ($1, $2, $3)`,
		productID, key, value,
	)
	if err != nil {
		log.Printf("Spec warning %s: %v", key, err)
	}
}

// ─── Categories ───────────────────────────────────────────────────────────────

func seedCategories() {
	log.Println("⏳ Seeding categories...")

	createCategory("Apple", "apple", "Apple products including iPhones and accessories", nil)
	createCategory("iPhones", "iphones", "Latest Apple iPhones", nil)
	createCategory("Apple Accessories", "apple-accessories", "Genuine Apple accessories", nil)
	createCategory("Samsung", "samsung", "Samsung phones and accessories", nil)
	createCategory("Samsung Phones", "samsung-phones", "Latest Samsung smartphones", nil)
	createCategory("Samsung Accessories", "samsung-accessories", "Genuine Samsung accessories", nil)

	log.Println("✅ Categories seeded")
}

// ─── Apple Phones ─────────────────────────────────────────────────────────────

func seedApplePhones() {
	log.Println("⏳ Seeding Apple iPhones...")

	var categoryID string
	db.QueryRow(`SELECT id FROM categories WHERE slug = 'iphones'`).Scan(&categoryID)

	phones := []struct {
		name     string
		slug     string
		variants []struct {
			storage string
			color   string
			price   float64
			sku     string
		}
		chip    string
		display string
		camera  string
	}{
		{
			name: "iPhone 17 Pro Max",
			slug: "iphone-17-pro-max",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"1TB", "Orange", 255000, "IP17PM-1TB-ORG"},
				{"1TB", "Blue", 255000, "IP17PM-1TB-BLU"},
				{"1TB", "Silver", 255000, "IP17PM-1TB-SLV"},
				{"512GB", "Orange", 217000, "IP17PM-512-ORG"},
				{"512GB", "Blue", 217000, "IP17PM-512-BLU"},
				{"512GB", "Silver", 223000, "IP17PM-512-SLV"},
				{"256GB", "Orange", 190000, "IP17PM-256-ORG"},
				{"256GB", "Blue", 190000, "IP17PM-256-BLU"},
				{"256GB", "Silver", 190000, "IP17PM-256-SLV"},
			},
			chip:    "A19 Pro",
			display: "6.9-inch Super Retina XDR",
			camera:  "48MP Fusion + 48MP Ultra Wide + 12MP Telephoto",
		},
		{
			name: "iPhone 17 Pro",
			slug: "iphone-17-pro",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"512GB", "Orange/Blue", 205000, "IP17P-512-ORGBLU"},
				{"256GB", "Orange/Blue", 170000, "IP17P-256-ORGBLU"},
				{"256GB", "Silver", 175000, "IP17P-256-SLV"},
			},
			chip:    "A19 Pro",
			display: "6.3-inch Super Retina XDR",
			camera:  "48MP Fusion + 48MP Ultra Wide + 12MP Telephoto",
		},
		{
			name: "iPhone 17",
			slug: "iphone-17",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"256GB", "All Colours", 117000, "IP17-256-ALL"},
			},
			chip:    "A19",
			display: "6.1-inch Super Retina XDR",
			camera:  "48MP Fusion + 12MP Ultra Wide",
		},
		{
			name: "iPhone 17 Air",
			slug: "iphone-17-air",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"256GB", "All Colours", 123000, "IP17AIR-256-ALL"},
			},
			chip:    "A19",
			display: "6.6-inch Super Retina XDR",
			camera:  "48MP Fusion + 12MP Ultra Wide",
		},
		{
			name: "iPhone 16 Pro Max",
			slug: "iphone-16-pro-max",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"1TB", "All Colours", 185000, "IP16PM-1TB-ALL"},
				{"512GB", "All Colours", 174000, "IP16PM-512-ALL"},
				{"256GB", "All Colours", 157000, "IP16PM-256-ALL"},
			},
			chip:    "A18 Pro",
			display: "6.9-inch Super Retina XDR",
			camera:  "48MP Fusion + 48MP Ultra Wide + 12MP Telephoto",
		},
		{
			name: "iPhone 16 Pro",
			slug: "iphone-16-pro",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"1TB", "All Colours", 200000, "IP16P-1TB-ALL"},
				{"512GB", "All Colours", 174000, "IP16P-512-ALL"},
				{"256GB", "All Colours", 146000, "IP16P-256-ALL"},
				{"128GB", "All Colours", 125000, "IP16P-128-ALL"},
			},
			chip:    "A18 Pro",
			display: "6.3-inch Super Retina XDR",
			camera:  "48MP Fusion + 48MP Ultra Wide + 12MP Telephoto",
		},
		{
			name: "iPhone 16 Plus",
			slug: "iphone-16-plus",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"256GB", "All Colours", 123000, "IP16PL-256-ALL"},
				{"128GB", "All Colours", 100000, "IP16PL-128-ALL"},
			},
			chip:    "A18",
			display: "6.7-inch Super Retina XDR",
			camera:  "48MP Fusion + 12MP Ultra Wide",
		},
		{
			name: "iPhone 16",
			slug: "iphone-16",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"256GB", "All Colours", 104000, "IP16-256-ALL"},
				{"128GB", "All Colours", 92000, "IP16-128-ALL"},
			},
			chip:    "A18",
			display: "6.1-inch Super Retina XDR",
			camera:  "48MP Fusion + 12MP Ultra Wide",
		},
		{
			name: "iPhone 16e",
			slug: "iphone-16e",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"256GB", "All Colours", 82000, "IP16E-256-ALL"},
				{"128GB", "All Colours", 70000, "IP16E-128-ALL"},
			},
			chip:    "A16",
			display: "6.1-inch Super Retina XDR",
			camera:  "48MP Fusion",
		},
		{
			name: "iPhone 15 Plus",
			slug: "iphone-15-plus",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"128GB", "All Colours", 92000, "IP15PL-128-ALL"},
			},
			chip:    "A16",
			display: "6.7-inch Super Retina XDR",
			camera:  "48MP Main + 12MP Ultra Wide",
		},
		{
			name: "iPhone 15",
			slug: "iphone-15",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"256GB", "All Colours", 99000, "IP15-256-ALL"},
				{"128GB", "All Colours", 79000, "IP15-128-ALL"},
			},
			chip:    "A16",
			display: "6.1-inch Super Retina XDR",
			camera:  "48MP Main + 12MP Ultra Wide",
		},
		{
			name: "iPhone 14",
			slug: "iphone-14",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"256GB", "All Colours", 82000, "IP14-256-ALL"},
				{"128GB", "All Colours", 71000, "IP14-128-ALL"},
			},
			chip:    "A15",
			display: "6.1-inch Super Retina XDR",
			camera:  "12MP Main + 12MP Ultra Wide",
		},
		{
			name: "iPhone 13",
			slug: "iphone-13",
			variants: []struct {
				storage string
				color   string
				price   float64
				sku     string
			}{
				{"128GB", "All Colours", 68000, "IP13-128-ALL"},
			},
			chip:    "A15",
			display: "6.1-inch Super Retina XDR",
			camera:  "12MP Main + 12MP Ultra Wide",
		},
	}

	for _, p := range phones {
		id := createProduct(
			p.name,
			p.slug,
			fmt.Sprintf("The %s features the %s chip, a stunning %s display and an advanced %s camera system.", p.name, p.chip, p.display, p.camera),
			categoryID,
			p.variants[0].price,
			fmt.Sprintf("Buy %s in Kenya | Best Price", p.name),
		)
		for _, v := range p.variants {
			createVariant(id, v.sku, v.color, v.storage, v.price, 10)
		}
		createSpec(id, "Chip", p.chip)
		createSpec(id, "Display", p.display)
		createSpec(id, "Camera", p.camera)
		createSpec(id, "Operating System", "iOS 18")
		log.Printf("   ✓ %s", p.name)
	}

	log.Println("✅ Apple iPhones seeded")
}

// ─── Apple Accessories ────────────────────────────────────────────────────────

func seedAppleAccessories() {
	log.Println("⏳ Seeding Apple Accessories...")

	var categoryID string
	db.QueryRow(`SELECT id FROM categories WHERE slug = 'apple-accessories'`).Scan(&categoryID)

	accessories := []struct {
		name  string
		slug  string
		price float64
		sku   string
		spec  string
	}{
		{"AirPods 4 ANC", "airpods-4-anc", 19500, "AP4ANC", "Active Noise Cancellation"},
		{"AirPods 4", "airpods-4", 14000, "AP4", "Personalised Spatial Audio"},
		{"AirPods Pro 3", "airpods-pro-3", 32000, "APP3", "Active Noise Cancellation + Transparency"},
		{"Apple Pencil Pro", "apple-pencil-pro", 16000, "APENCPRO", "Squeeze gesture, Find My"},
		{"Apple Pencil 2nd Gen", "apple-pencil-2nd-gen", 12000, "APENC2", "Double-tap gesture"},
		{"Apple Pencil USB-C", "apple-pencil-usbc", 12000, "APENCUSBC", "USB-C connectivity"},
		{"Apple Pencil 1st Gen", "apple-pencil-1st-gen", 12000, "APENC1", "Lightning connectivity"},
		{"MagSafe Charger", "magsafe-charger", 7000, "MAGSAFE", "15W MagSafe wireless charging"},
		{"Apple Watch Charger USB-C", "apple-watch-charger-usbc", 4500, "AWCHARGER", "USB-C magnetic charging"},
		{"AirTag 1 Pack", "airtag-1-pack", 4500, "ATAG1", "Precision Finding with Ultra Wideband"},
		{"AirTag 4 Pack", "airtag-4-pack", 12000, "ATAG4", "Precision Finding with Ultra Wideband"},
		{"Magic Mouse 3 Black", "magic-mouse-3-black", 13500, "MM3BLK", "Multi-Touch surface"},
		{"Magic Mouse 3 Silver", "magic-mouse-3-silver", 10500, "MM3SLV", "Multi-Touch surface"},
		{"Magic Trackpad Type-C", "magic-trackpad-typec", 17500, "MTTPC", "Force Touch, USB-C"},
		{"Magic Trackpad Lightning", "magic-trackpad-lightning", 17000, "MTTLT", "Force Touch, Lightning"},
		{"240W Cable USB-C Mac", "240w-cable-usbc", 3500, "CBL240W", "240W charging support"},
		{"Apple TV 4K 128GB", "apple-tv-4k-128gb", 24000, "ATV4K128", "4K HDR, Dolby Vision"},
		{"Apple Thunderbolt 4 Pro Cable 1.8m", "thunderbolt-4-pro-cable", 108000, "TB4PRO", "40Gb/s data transfer"},
		{"35W Adapter USB-C Mac", "35w-adapter-usbc", 7500, "ADP35W", "Dual USB-C ports"},
		{"70W Adapter USB-C Mac", "70w-adapter-usbc", 9500, "ADP70W", "USB-C"},
		{"96W Adapter USB-C Mac", "96w-adapter-usbc", 11500, "ADP96W", "USB-C"},
		{"140W Adapter USB-C Mac", "140w-adapter-usbc", 13500, "ADP140W", "USB-C"},
		{"Apple USB-C MagSafe 3 Cable", "usbc-magsafe-3-cable", 7000, "CBLMAGSAFE3", "MagSafe 3 connector"},
		{"iPhone Air MagSafe Battery", "iphone-air-magsafe-battery", 16000, "MAGBAT", "MagSafe wireless charging"},
		{"Apple Pencil Tips", "apple-pencil-tips", 4000, "IPTIPS", "4 replacement tips"},
		{"Apple Watch Ultra 3 49mm", "apple-watch-ultra-3-49mm", 99000, "AWU349", "Titanium, 49mm"},
		{"Apple Watch Ultra 3 49mm Milanese", "apple-watch-ultra-3-49mm-milanese", 123000, "AWU349MIL", "Titanium, Milanese Loop"},
		{"Apple Watch Series 11 46mm", "apple-watch-series-11-46mm", 50000, "AWS1146", "Always-On Retina display"},
		{"Apple Watch Series 11 42mm", "apple-watch-series-11-42mm", 46000, "AWS1142", "Always-On Retina display"},
		{"Apple Watch Series 10 46mm Rose Gold", "apple-watch-series-10-46mm", 39000, "AWS1046RG", "Rose Gold"},
		{"Apple Watch SE 3 40mm Black", "apple-watch-se-3-40mm-black", 34000, "AWSE340BLK", "40mm Black"},
		{"Apple Watch SE 3 40mm Starlight", "apple-watch-se-3-40mm-starlight", 34000, "AWSE340STR", "40mm Starlight"},
		{"Apple Watch SE 3 44mm", "apple-watch-se-3-44mm", 37000, "AWSE344", "44mm"},
		{"Magic Keyboard 13 M4 M5", "magic-keyboard-13-m4-m5", 54000, "MKB13M4", "Touch ID, USB-C"},
		{"Magic Keyboard 11 M4 M5", "magic-keyboard-11-m4-m5", 43000, "MKB11M4", "Touch ID, USB-C"},
		{"Magic Keyboard 10/11 Gen", "magic-keyboard-10-11-gen", 34000, "MKB1011", "Touch ID"},
		{"Magic Keyboard Air 13 M3 Silver", "magic-keyboard-air-13-m3", 38000, "MKBAIR13", "Touch ID, USB-C"},
		{"Magic Keyboard 11 M1 M2", "magic-keyboard-11-m1-m2", 34000, "MKB11M1", "Touch ID"},
		{"Magic Keyboard 12.9 M1 M2", "magic-keyboard-12-9-m1-m2", 32000, "MKB12M1", "Touch ID"},
		{"Folio Keyboard 11", "folio-keyboard-11", 24000, "FKBD11", "Smart Connector"},
		{"Magic Keyboard Touch ID Numeric", "magic-keyboard-touch-id-numeric", 27000, "MKBTIDNUM", "Touch ID, Numeric keypad"},
		{"Magic Keyboard Touch ID Non Numeric", "magic-keyboard-touch-id-non-numeric", 20000, "MKBTIDNON", "Touch ID"},
		{"EarPods Lightning", "earpods-lightning", 3000, "EPL", "Lightning connector"},
		{"EarPods Type-C", "earpods-typec", 3000, "EPTC", "USB-C connector"},
		{"20W Adapter", "apple-20w-adapter", 3000, "ADP20W", "USB-C Power Adapter"},
		{"35W Dual USB-C Power Adapter", "apple-35w-dual-usbc", 7500, "ADP35WDUAL", "Two USB-C ports"},
		{"Lightning to iPhone Jack", "lightning-to-iphone-jack", 1500, "LTIJ", "3.5mm headphone jack adapter"},
		{"Lightning to Type-C", "lightning-to-typec", 3000, "LTTC", "Lightning to USB-C cable"},
		{"C to C Cable", "c-to-c-cable", 3000, "CTOC", "USB-C to USB-C cable"},
		{"World Travel Kit Adapter", "world-travel-kit-adapter", 6000, "WTK", "International plug adapters"},
	}

	for _, a := range accessories {
		id := createProduct(
			a.name,
			a.slug,
			fmt.Sprintf("Genuine Apple %s. %s.", a.name, a.spec),
			categoryID,
			a.price,
			fmt.Sprintf("Buy %s in Kenya | Best Price", a.name),
		)
		createVariant(id, a.sku, "Default", "N/A", a.price, 20)
		createSpec(id, "Features", a.spec)
		createSpec(id, "Brand", "Apple")
		log.Printf("   ✓ %s", a.name)
	}

	log.Println("✅ Apple Accessories seeded")
}

// ─── Samsung Phones ───────────────────────────────────────────────────────────

func seedSamsungPhones() {
	log.Println("⏳ Seeding Samsung Phones...")

	var categoryID string
	db.QueryRow(`SELECT id FROM categories WHERE slug = 'samsung-phones'`).Scan(&categoryID)

	phones := []struct {
		name     string
		slug     string
		variants []struct {
			storage string
			price   float64
			sku     string
		}
		display string
		camera  string
	}{
		{
			name: "Samsung Galaxy Fold 7",
			slug: "samsung-galaxy-fold-7",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"512GB", 198000, "SGF7-512"},
				{"256GB", 181000, "SGF7-256"},
			},
			display: "7.6-inch Foldable Dynamic AMOLED 2X",
			camera:  "200MP Main + 12MP Ultra Wide + 10MP Telephoto",
		},
		{
			name: "Samsung Galaxy S26 Ultra",
			slug: "samsung-galaxy-s26-ultra",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"1TB", 198000, "SGS26U-1TB"},
				{"512GB", 149000, "SGS26U-512"},
				{"256GB", 134000, "SGS26U-256"},
			},
			display: "6.9-inch Dynamic AMOLED 2X",
			camera:  "200MP Main + 12MP Ultra Wide + 50MP Telephoto",
		},
		{
			name: "Samsung Galaxy S26 Plus",
			slug: "samsung-galaxy-s26-plus",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"256GB", 118000, "SGS26PL-256"},
			},
			display: "6.7-inch Dynamic AMOLED 2X",
			camera:  "50MP Main + 12MP Ultra Wide + 10MP Telephoto",
		},
		{
			name: "Samsung Galaxy S26",
			slug: "samsung-galaxy-s26",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"256GB", 106000, "SGS26-256"},
			},
			display: "6.2-inch Dynamic AMOLED 2X",
			camera:  "50MP Main + 12MP Ultra Wide + 10MP Telephoto",
		},
		{
			name: "Samsung Galaxy S25 Ultra",
			slug: "samsung-galaxy-s25-ultra",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"512GB", 136000, "SGS25U-512"},
				{"256GB", 117000, "SGS25U-256"},
			},
			display: "6.9-inch Dynamic AMOLED 2X",
			camera:  "200MP Main + 12MP Ultra Wide + 50MP Telephoto",
		},
		{
			name: "Samsung Galaxy S25 FE",
			slug: "samsung-galaxy-s25-fe",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"256GB", 72000, "SGS25FE-256"},
			},
			display: "6.7-inch Dynamic AMOLED 2X",
			camera:  "50MP Main + 12MP Ultra Wide",
		},
		{
			name: "Samsung Galaxy S25",
			slug: "samsung-galaxy-s25",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"256GB", 93000, "SGS25-256"},
			},
			display: "6.2-inch Dynamic AMOLED 2X",
			camera:  "50MP Main + 12MP Ultra Wide + 10MP Telephoto",
		},
		{
			name: "Samsung Galaxy Flip 7",
			slug: "samsung-galaxy-flip-7",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"256GB", 105000, "SGFLIP7-256"},
			},
			display: "6.7-inch Foldable Dynamic AMOLED 2X",
			camera:  "50MP Main + 12MP Ultra Wide",
		},
		{
			name: "Samsung Galaxy Flip 7 FE",
			slug: "samsung-galaxy-flip-7-fe",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"256GB", 91000, "SGFLIP7FE-256"},
				{"128GB", 81000, "SGFLIP7FE-128"},
			},
			display: "6.7-inch Foldable Dynamic AMOLED",
			camera:  "50MP Main + 12MP Ultra Wide",
		},
		{
			name: "Samsung Galaxy A57",
			slug: "samsung-galaxy-a57",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"256GB", 55000, "SGA57-256"},
			},
			display: "6.7-inch Super AMOLED",
			camera:  "64MP Main + 12MP Ultra Wide",
		},
		{
			name: "Samsung Galaxy A37",
			slug: "samsung-galaxy-a37",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"256GB", 48000, "SGA37-256"},
			},
			display: "6.5-inch Super AMOLED",
			camera:  "50MP Main + 8MP Ultra Wide",
		},
		{
			name: "Samsung Galaxy A17",
			slug: "samsung-galaxy-a17",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"256GB", 24800, "SGA17-8-256"},
				{"128GB", 20600, "SGA17-6-128"},
				{"128GB", 18800, "SGA17-4-128"},
			},
			display: "6.5-inch PLS LCD",
			camera:  "50MP Main + 5MP Ultra Wide",
		},
		{
			name: "Samsung Galaxy A16",
			slug: "samsung-galaxy-a16",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"256GB", 22000, "SGA16-8-256"},
				{"128GB", 17500, "SGA16-6-128"},
				{"128GB", 16000, "SGA16-4-128"},
			},
			display: "6.7-inch Super AMOLED",
			camera:  "50MP Main + 5MP Ultra Wide",
		},
		{
			name: "Samsung Galaxy A07",
			slug: "samsung-galaxy-a07",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"128GB", 13500, "SGA07-4-128"},
				{"64GB", 12500, "SGA07-4-64"},
			},
			display: "6.5-inch PLS LCD",
			camera:  "50MP Main",
		},
		{
			name: "Samsung Galaxy A06",
			slug: "samsung-galaxy-a06",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"128GB", 12500, "SGA06-6-128"},
				{"128GB", 12000, "SGA06-4-128"},
				{"64GB", 11500, "SGA06-4-64"},
			},
			display: "6.7-inch PLS LCD",
			camera:  "50MP Main",
		},
		{
			name: "Samsung Galaxy A05s",
			slug: "samsung-galaxy-a05s",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"128GB", 13500, "SGA05S-4-128"},
			},
			display: "6.7-inch PLS LCD",
			camera:  "50MP Main + 2MP Depth",
		},
		{
			name: "Samsung Galaxy A05",
			slug: "samsung-galaxy-a05",
			variants: []struct {
				storage string
				price   float64
				sku     string
			}{
				{"64GB", 10000, "SGA05-4-64"},
			},
			display: "6.7-inch PLS LCD",
			camera:  "50MP Main",
		},
	}

	for _, p := range phones {
		id := createProduct(
			p.name,
			p.slug,
			fmt.Sprintf("The %s features a stunning %s display and a powerful %s camera system.", p.name, p.display, p.camera),
			categoryID,
			p.variants[0].price,
			fmt.Sprintf("Buy %s in Kenya | Best Price", p.name),
		)
		for _, v := range p.variants {
			createVariant(id, v.sku, "All Colours", v.storage, v.price, 10)
		}
		createSpec(id, "Display", p.display)
		createSpec(id, "Camera", p.camera)
		createSpec(id, "Operating System", "Android 15")
		createSpec(id, "Brand", "Samsung")
		log.Printf("   ✓ %s", p.name)
	}

	log.Println("✅ Samsung Phones seeded")
}

// ─── Samsung Accessories ──────────────────────────────────────────────────────

func seedSamsungAccessories() {
	log.Println("⏳ Seeding Samsung Accessories...")

	var categoryID string
	db.QueryRow(`SELECT id FROM categories WHERE slug = 'samsung-accessories'`).Scan(&categoryID)

	accessories := []struct {
		name  string
		slug  string
		price float64
		sku   string
		spec  string
	}{
		{"Samsung Galaxy Buds 4 Pro", "samsung-buds-4-pro", 25500, "SGB4PRO", "Active Noise Cancellation, 360 Audio"},
		{"Samsung Galaxy Buds 4", "samsung-buds-4", 17500, "SGB4", "Active Noise Cancellation"},
		{"Samsung Galaxy Buds 3 FE", "samsung-buds-3-fe", 10000, "SGB3FE", "ANC Lite"},
		{"Samsung Galaxy Buds Core", "samsung-buds-core", 4000, "SGBCORE", "Up to 30hrs battery"},
		{"Samsung Galaxy Buds 3", "samsung-buds-3", 9000, "SGB3", "Blade-type design"},
		{"Samsung Galaxy Watch Ultra 2025", "samsung-galaxy-watch-ultra-2025", 43000, "SGWU2025", "47mm Titanium"},
		{"Samsung Galaxy Watch 8 Classic", "samsung-galaxy-watch-8-classic", 34000, "SGW8C", "47mm Stainless Steel"},
		{"Samsung Galaxy Watch 8 44mm", "samsung-galaxy-watch-8-44mm", 28500, "SGW844", "Aluminum"},
		{"Samsung Galaxy Watch 8 40mm", "samsung-galaxy-watch-8-40mm", 27000, "SGW840", "Aluminum"},
		{"Samsung Galaxy Watch FE", "samsung-galaxy-watch-fe", 12500, "SGWFE", "40mm Aluminum"},
		{"Samsung Smart Tag 2 1 Pack", "samsung-smart-tag-2-1pack", 2500, "SGTAG21", "UWB Precision Finding"},
		{"Samsung Smart Tag 2 4 Pack", "samsung-smart-tag-2-4pack", 7500, "SGTAG24", "UWB Precision Finding"},
		{"Samsung 45W Adapter with 2m Cable", "samsung-45w-adapter-cable", 3500, "SADP45WC", "Super Fast Charging 2.0"},
		{"Samsung 45W Adapter", "samsung-45w-adapter", 2800, "SADP45W", "Super Fast Charging 2.0"},
		{"Samsung 25W Adapter with Cable", "samsung-25w-adapter-cable", 2500, "SADP25WC", "Fast Charging"},
		{"Samsung 25W Adapter", "samsung-25w-adapter", 1500, "SADP25W", "Fast Charging"},
		{"Samsung Wireless Charger Duo", "samsung-wireless-charger-duo", 7500, "SWCDUO", "Charges phone and watch"},
		{"Samsung Wireless Charger Trio", "samsung-wireless-charger-trio", 8500, "SWCTRIO", "Charges phone, watch and buds"},
		{"Samsung Tab S11 12/256GB", "samsung-tab-s11", 99000, "SGT S11-12-256", "12.2-inch Dynamic AMOLED 2X"},
		{"Samsung Tab S10 FE Plus 12/256GB 5G", "samsung-tab-s10-fe-plus", 83000, "SGTS10FEPLUS", "12.4-inch TFT LCD"},
		{"Samsung Tab S10 FE 256GB 5G", "samsung-tab-s10-fe-256-5g", 70000, "SGTS10FE256-5G", "10.9-inch TFT LCD"},
		{"Samsung Tab S10 FE 128GB 5G", "samsung-tab-s10-fe-128-5g", 55000, "SGTS10FE128-5G", "10.9-inch TFT LCD"},
		{"Samsung Tab S10 FE 256GB WiFi", "samsung-tab-s10-fe-256-wifi", 53000, "SGTS10FE256-WIFI", "10.9-inch TFT LCD"},
		{"Samsung Tab S10 FE 128GB WiFi", "samsung-tab-s10-fe-128-wifi", 50000, "SGTS10FE128-WIFI", "10.9-inch TFT LCD"},
		{"Samsung Tab S10 Lite 8/256GB 5G", "samsung-tab-s10-lite-8-256-5g", 60000, "SGTS10L-8-256-5G", "10.1-inch TFT LCD"},
		{"Samsung Tab S10 Lite 8/256GB WiFi", "samsung-tab-s10-lite-8-256-wifi", 46000, "SGTS10L-8-256-WIFI", "10.1-inch TFT LCD"},
		{"Samsung Tab S10 Lite 6/128GB WiFi", "samsung-tab-s10-lite-6-128-wifi", 43000, "SGTS10L-6-128-WIFI", "10.1-inch TFT LCD"},
		{"Samsung Tab S6 Lite 4/64GB WiFi Dubai", "samsung-tab-s6-lite", 30000, "SGTS6L-4-64", "10.4-inch TFT LCD"},
		{"Samsung Tab A11 4/64GB", "samsung-tab-a11", 15800, "SGTA11-4-64", "10.4-inch TFT LCD"},
		{"Samsung S9/S9FE/S10FE/S9 5G Keyboard", "samsung-keyboard-s9", 21000, "SGKBD-S9", "Book Cover Keyboard"},
		{"Samsung S10 FE+/S10 FE+5G Keyboard", "samsung-keyboard-s10fe", 25000, "SGKBD-S10FE", "Book Cover Keyboard"},
		{"Samsung S11 Ultra Keyboard", "samsung-keyboard-s11-ultra", 25000, "SGKBD-S11U", "Book Cover Keyboard"},
		{"Samsung S10 Ultra Keyboard", "samsung-keyboard-s10-ultra", 37000, "SGKBD-S10U", "Book Cover Keyboard"},
		{"Samsung S10 Plus/S9FE+/S9 Plus Keyboard", "samsung-keyboard-s10-plus", 25000, "SGKBD-S10PL", "Book Cover Keyboard"},
		{"Samsung S11/S11 5G Keyboard", "samsung-keyboard-s11", 22000, "SGKBD-S11", "Book Cover Keyboard"},
		{"Samsung Tab S10 FE Keyboard", "samsung-keyboard-tab-s10fe", 21000, "SGKBD-T10FE", "Book Cover Keyboard"},
		{"Samsung Tab A9 Plus Keyboard", "samsung-keyboard-tab-a9plus", 10500, "SGKBD-TA9PL", "Book Cover Keyboard"},
	}

	for _, a := range accessories {
		id := createProduct(
			a.name,
			a.slug,
			fmt.Sprintf("Genuine Samsung %s. %s.", a.name, a.spec),
			categoryID,
			a.price,
			fmt.Sprintf("Buy %s in Kenya | Best Price", a.name),
		)
		createVariant(id, a.sku, "Default", "N/A", a.price, 20)
		createSpec(id, "Features", a.spec)
		createSpec(id, "Brand", "Samsung")
		log.Printf("   ✓ %s", a.name)
	}

	log.Println("✅ Samsung Accessories seeded")
}

// ─── Admin User ───────────────────────────────────────────────────────────────

func seedAdminUser() {
	log.Println("⏳ Seeding admin user...")

	var count int
	db.QueryRow(`SELECT COUNT(*) FROM users WHERE email = 'admin@applestore.co.ke'`).Scan(&count)
	if count > 0 {
		log.Println("   ℹ Admin user already exists, skipping")
		return
	}

	// bcrypt hash of "Admin@1234"
	hash := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"

	_, err := db.Exec(`
		INSERT INTO users (name, email, password_hash, provider, is_verified, is_active, role)
		VALUES ('Admin', 'admin@applestore.co.ke', $1, 'local', TRUE, TRUE, 'admin')`,
		hash,
	)
	if err != nil {
		log.Printf("Admin user warning: %v", err)
		return
	}

	log.Println("   ✓ Admin user created")
	log.Println("   📧 Email: admin@applestore.co.ke")
	log.Println("   🔑 Password: Admin@1234")
	log.Println("✅ Admin user seeded")
}
