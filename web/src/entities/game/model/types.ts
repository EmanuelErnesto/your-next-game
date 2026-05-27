export interface Game {
  id: string;
  title: string;
  coverUrl: string;
  status: 'Jogado' | 'Não Jogado';
  genre: string;
  platform: string;
  description: string;
  developer: string;
  releaseDate: string;
  hoursPlayed: number;
  rating: number | null;
  steamAppId: string | null;
}
