
importScripts('https://www.gstatic.com/firebasejs/8.10.1/firebase-app.js');
importScripts('https://www.gstatic.com/firebasejs/8.10.1/firebase-messaging.js');

const firebaseConfig = {
  apiKey: "AIzaSyA9_JMiFr8nU6l8AYYXwBKgFRCl4_O39O0",
  authDomain: "subash-bakery.firebaseapp.com",
  projectId: "subash-bakery",
  storageBucket: "subash-bakery.firebasestorage.app",
  messagingSenderId: "470400597578",
  appId: "1:470400597578:web:e5588dd32efc48403ff0b3",
};

firebase.initializeApp(firebaseConfig);
const messaging = firebase.messaging();

messaging.onBackgroundMessage((payload) => {
  console.log('[firebase-messaging-sw.js] Received background message ', payload);
  const notificationTitle = payload.notification.title;
  const notificationOptions = {
    body: payload.notification.body,
    icon: '/vite.svg'
  };

  self.registration.showNotification(notificationTitle, notificationOptions);
});
