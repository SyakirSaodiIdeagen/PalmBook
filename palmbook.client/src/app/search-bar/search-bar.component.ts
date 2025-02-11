import { Component, ElementRef, ViewChild } from '@angular/core';

@Component({
  selector: 'app-search-bar',
  standalone: false,
  templateUrl: './search-bar.component.html',
  styleUrl: './search-bar.component.css'
})
export class SearchBarComponent {
  query: string = '';
  results: { name: string; type: string; source: string; downloadUrl: string }[] = [];
  loading = false;
  searched = false;

  @ViewChild('topElement') topElement!: ElementRef;

  search() {
    if (!this.query.trim()) return;

    this.loading = true;
    this.results = [];
    this.searched = true;

    setTimeout(() => {
      this.loading = false;
      this.results =
        this.query.toLowerCase() === 'test'
          ? [
            { name: 'Report.pdf', type: 'PDF', source: 'SharePoint', downloadUrl: '#' },
            { name: 'Data.xlsx', type: 'Excel', source: 'Teams', downloadUrl: '#' },
            { name: 'Presentation.pptx', type: 'PowerPoint', source: 'OneDrive', downloadUrl: '#' }
          ]
          : [];

      this.scrollToTop();
    }, 1500);
  }

  scrollToTop() {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }
}
