import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

function Notifications() {
    const { user_id } = useParams();
    const [notifications, setNotifications] = useState([]);

    useEffect(() => {
        const ws = new WebSocket(`ws://localhost:8080/notifications/${user_id}`);

        ws.onmessage = (event) => {
            const newNotification = JSON.parse(event.data);
            setNotifications((prevNotifications) => [newNotification, ...prevNotifications]);
        };

        return () => {
            ws.close();
        };
    }, [user_id]);

    return (
        <div>
            <h1>Notifications</h1>
            <ul>
                {notifications.map(notification => (
                    <li key={notification.id}>
                        {notification.type} - {notification.post_id}
                    </li>
                ))}
            </ul>
        </div>
    );
}

export default Notifications;
