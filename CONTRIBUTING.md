Branching & PR Constraint
- Repo bisa diakses semua member
- Branch utama: `main` dan `dev`
- Branch rules: fitur dibuat di feature/<nama fitur>
- Ada file CONTRIBUTING.md singkat berisi aturan commit & Pull Request (PR)

Branching Rules:
- Penamaan branch harus kategori/nama. misal: feature/<nama fitur>
- Commit message harus menyertakan kode task, misal “[KAN-13] membuat tombol”
- Untuk Pull Request, source branch: dev
- Setelah merge ke main, hapus branch yang sudah ga dipake


PR Rules:
<!-- Judul PR -->
[KAN-13] Udah nambahin fitur A

🔗 Related Jira Issue
- Jira: [KAN-13] Create Repo & Branching Strategy



🎯 Goal / What problem does this PR solve?
<!-- Jelaskan singkat tujuan PR -->
Add feature to create a food entry and update daily total calories.



✅ What changed?
<!-- Checklist perubahan yang dilakukan -->
- [✅] Added new feature A.
- [ ] Updated existing logic: ...
- [ ] Refactored code: ...
- [ ] Fixed bug: ...



🧩 Implementation Details
<!-- Jelaskan cara kerja, bukan cuma "sudah dibuat" -->
- Data flow:
  - Input: Membuat button untuk cek kalori
  - Process: Menggunakan function blablala dan menghubungkan blablabla untuk logicnya
  - Output: Bisa cek kalori
- Key files changed:
  - `src/feature/kalori`
  - `src/...`



🧪 How to test (Step-by-step)
<!-- WAJIB: reviewer harus bisa test tanpa nebak -->
1. Run app: `npm run dev`
2. Go to page: `localhost:3000/cekkalori`
3. Input:
   - Pencet tombol cek kalori makanan
4. Expected result:
   - Muncul kalori dan nutrisi dari makanan



📸 Screenshots / UI Preview (if UI related)
<!-- Optional tapi disarankan -->
Before:
- (attach screenshot)

After:
- (attach screenshot)



🧠 Learning Notes (Mandatory for this project)
<!-- Ini bagian skill development: tulis 2–5 poin -->
What I learned:
- ...
- ...

What was difficult / still confusing:
- ...
- ...



⚠️ Risks / Edge Cases
<!-- Misal: data kosong, input invalid, crash scenario -->
- [ ] Handles empty data
- [ ] Handles invalid input
- [ ] No crash when user profile not set

Notes:
- ...



🔍 Checklist (Definition of Done)
Code sudah di-push ke branch feature

- [ ] Ada PR ke dev

- [ ] Minimal 1 reviewer approve (experienced/intermediate)

- [ ] Tidak ada error runtime di flow yang dikerjakan

- [ ] Acceptance criteria issue terpenuhi

- [ ] Issue status = Done



📌 Reviewer Notes
<!-- Apa yang reviewer harus fokus cek? -->
Please review:
- [ ] Logic correctness
- [✅ ] UI/UX behavior
- [✅ ] Code structure
- [ ] Potensi bug/error


