<?php

namespace App\EventSubscriber;

use Symfony\Component\EventDispatcher\EventSubscriberInterface;
use Symfony\Component\HttpFoundation\IpUtils;
use Symfony\Component\HttpFoundation\JsonResponse;
use Symfony\Component\HttpFoundation\Request;
use Symfony\Component\HttpKernel\Event\RequestEvent;
use Symfony\Component\HttpKernel\KernelEvents;

/**
 * Optional access gate for the API.
 *
 * Off unless SR_API_KEYS is set (comma-separated, so keys can be rotated or given per client). When it
 * is set, every /api/* request except /api/ping must either send one of the keys in the X-SR-Key header,
 * or come from an address in SR_ALLOWED_IPS (comma-separated IPs or CIDR ranges, optional). "/" and
 * /api/ping stay open for health checks.
 *
 * The client address is the right-most public address in X-Forwarded-For: the platform's proxies append
 * the real address there, and a caller can only add entries to the left of it.
 */
class ApiAccessSubscriber implements EventSubscriberInterface
{
    private const OPEN_PATHS = ['/api/ping'];

    public static function getSubscribedEvents(): array
    {
        // Before routing (32), so no controller work happens for a refused call.
        return [KernelEvents::REQUEST => ['onKernelRequest', 256]];
    }

    private static function setting(string $name): string
    {
        $value = $_SERVER[$name] ?? $_ENV[$name] ?? getenv($name);
        return is_string($value) ? trim($value) : '';
    }

    private static function list(string $value): array
    {
        return array_values(array_filter(array_map('trim', explode(',', $value)), 'strlen'));
    }

    public static function clientIp(Request $request): string
    {
        $chain = self::list((string) $request->headers->get('X-Forwarded-For', ''));
        for ($i = count($chain) - 1; $i >= 0; $i--) {
            $ip = $chain[$i];
            if (filter_var($ip, FILTER_VALIDATE_IP, FILTER_FLAG_NO_PRIV_RANGE | FILTER_FLAG_NO_RES_RANGE)) {
                return $ip;
            }
        }
        return (string) $request->server->get('REMOTE_ADDR', '');
    }

    public function onKernelRequest(RequestEvent $event): void
    {
        if (!$event->isMainRequest()) {
            return;
        }
        $request = $event->getRequest();
        $path = $request->getPathInfo();
        if (strpos($path, '/api/') !== 0 && $path !== '/api') {
            return;
        }
        if (in_array(rtrim($path, '/'), self::OPEN_PATHS, true)) {
            return;
        }
        $keys = self::list(self::setting('SR_API_KEYS'));
        if (!$keys) {
            return;                                     // gate not enabled
        }

        $sent = (string) $request->headers->get('X-SR-Key', '');
        if ($sent !== '') {
            foreach ($keys as $key) {
                if (hash_equals($key, $sent)) {
                    return;
                }
            }
        }
        $allowed = self::list(self::setting('SR_ALLOWED_IPS'));
        if ($allowed && IpUtils::checkIp(self::clientIp($request), $allowed)) {
            return;
        }

        $event->setResponse(new JsonResponse(
            ['status' => 401, 'message' => 'This API requires an access key (X-SR-Key header).'],
            401
        ));
    }
}
