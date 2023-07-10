export interface Notification {
  id: string;
  severity: 'error' | 'warning';
  message: React.ReactNode;
}
