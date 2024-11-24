// ...kodenya yang ada...

function setCookie(name, value, days) {
    const expires = new Date();
    expires.setTime(expires.getTime() + days * 24 * 60 * 60 * 1000);
    document.cookie = `${name}=${value};expires=${expires.toUTCString()};path=/;secure;samesite=strict`;
}

function getCookie(name) {
    const value = `; ${document.cookie}`;
    const parts = value.split(`; ${name}=`);
    if (parts.length === 2) return parts.pop().split(';').shift();
}

document.getElementById('register-form').addEventListener('submit', function(e) {
    e.preventDefault();
    const data = {
        username: this.username.value,
        email: this.email.value,
        password: this.password.value
    };
    fetch('/register', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            alert(data.message);
        }
    })
    .catch(error => console.error('Error:', error));
});

document.getElementById('login-form').addEventListener('submit', function(e) {
    e.preventDefault();
    const data = {
        email: this.email.value,
        password: this.password.value
    };
    fetch('/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            // Store tokens in cookies
            setCookie('access_token', data.access_token, 1); // 1 day
            setCookie('refresh_token', data.refresh_token, 7); // 7 days
            
            alert('Login successful\nAccess Token: ' + data.access_token + 
                  '\nRefresh Token: ' + data.refresh_token);
        }
    })
    .catch(error => console.error('Error:', error));
});

document.getElementById('change-password-form').addEventListener('submit', function(e) {
    e.preventDefault();
    const data = {
        old_password: this.old_password.value,
        new_password: this.new_password.value
    };
    
    const accessToken = getCookie('access_token');
    
    fetch('/auth/change-password', {
        method: 'POST',
        headers: { 
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${accessToken}`
        },
        body: JSON.stringify(data)
    })
    .then(response => {
        if (response.status === 401) {
            // Token expired, try to refresh
            return refreshToken().then(() => {
                // Retry with new access token
                const newAccessToken = getCookie('access_token');
                return fetch('/auth/change-password', {
                    method: 'POST',
                    headers: { 
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${newAccessToken}`
                    },
                    body: JSON.stringify(data)
                });
            });
        }
        return response.json();
    })
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            alert(data.message);
        }
    })
    .catch(error => console.error('Error:', error));
});

function refreshToken() {
    const refreshToken = getCookie('refresh_token');
    
    return fetch('/refresh-token', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${refreshToken}`
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.access_token) {
            setCookie('access_token', data.access_token, 1);
            return data;
        } else {
            throw new Error('Failed to refresh token');
        }
    });
}