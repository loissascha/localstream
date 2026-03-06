export function setCookie(name: string, value: string, days = 7) {
    const maxAge = days * 24 * 60 * 60;
    document.cookie = `${name}=${encodeURIComponent(value)}; path=/; max-age=${maxAge}`;
}

export function getCookie(name: string): string | null {
    const cookies = document.cookie.split('; ');

    for (const c of cookies) {
        const [key, value] = c.split('=');
        if (key === name) return decodeURIComponent(value);
    }

    return null;
}

export function deleteCookie(name: string) {
    document.cookie = `${name}=; path=/; max-age=0`;
}
