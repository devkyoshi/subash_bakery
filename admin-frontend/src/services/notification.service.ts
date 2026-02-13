import { messaging } from "../lib/firebase";
import { getToken, onMessage } from "firebase/messaging";
import axios from "axios";
import type { Notification as AppNotification } from "../types/notification.types";

const VAPID_KEY = import.meta.env.VITE_FIREBASE_VAPID_KEY; // Optional, if using VAPID
const API_URL = import.meta.env.VITE_API_BASE_URL || "/api/v1";

export const requestForToken = async () => {
  try {
    const permission = Notification.permission;
    if (permission === "denied") {
      console.warn("Notification permission has been denied by the user.");
      return null;
    }

    const currentToken = await getToken(messaging, { vapidKey: VAPID_KEY });
    if (currentToken) {
      // Send the token to your server
      await registerDevice(currentToken);
      return currentToken;
    } else {
      console.log(
        "No registration token available. Request permission to generate one.",
      );
      return null;
    }
  } catch (err) {
    console.log("An error occurred while retrieving token. ", err);
    return null;
  }
};

export const onMessageListener = (callback: (payload: any) => void) => {
  return onMessage(messaging, (payload) => {
    callback(payload);
  });
};

const registerDevice = async (token: string) => {
  try {
    const response = await axios.post(
      `${API_URL}/devices`,
      {
        token,
        platform: "web",
        name: navigator.userAgent,
      },
      {
        headers: {
          // Assumes auth token is handled by interceptor or we need to get it
          Authorization: `Bearer ${localStorage.getItem("accessToken")}`,
        },
      },
    );
  } catch (error) {
    console.error("Error registering device:", error);
  }
};

export const getNotifications = async (): Promise<AppNotification[]> => {
  try {
    const response = await axios.get(`${API_URL}/notifications`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("accessToken")}`,
      },
    });
    return response.data.data || [];
  } catch (error) {
    console.error("Failed to fetch notifications:", error);
    return [];
  }
};

export const markAsRead = async (id: string): Promise<void> => {
  try {
    await axios.patch(
      `${API_URL}/notifications/${id}/read`,
      {},
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("accessToken")}`,
        },
      },
    );
  } catch (error) {
    console.error("Failed to mark notification as read:", error);
  }
};

export const markAllAsRead = async (): Promise<void> => {
  try {
    await axios.patch(
      `${API_URL}/notifications/read-all`,
      {},
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem("accessToken")}`,
        },
      },
    );
  } catch (error) {
    console.error("Failed to mark all notifications as read:", error);
  }
};
