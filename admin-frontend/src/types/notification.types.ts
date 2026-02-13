export type NotificationType = "info" | "success" | "warning" | "error";

export interface Notification {
  id: string;
  title: string;
  message: string;
  read: boolean;
  type: NotificationType;
  createdAt: string; // ISO date string
  link?: string;
}
