<?php
/**
 * @author    Saman Esmaeil
 * @link      http://github.com/teal33t/jyotish-api for the canonical source repository
 * @license   GNU General Public License version 3 or later
 */


namespace Jyotish;

use Jyotish\Base\Data;
use Jyotish\Base\Locality;
use Jyotish\Base\Analysis;
use Jyotish\Ganita\Method\Swetest;
use Jyotish\Dasha\Dasha;
use Jyotish\Panchanga\AngaDefiner;
use Jyotish\Graha\Lagna;
use Jyotish\Yoga\Yoga;
use Jyotish\Bala\AshtakaVarga;
use Jyotish\Bala\GrahaBala;
use Jyotish\Bala\RashiBala;
use Jyotish\Graha\Graha;
use Carbon\Carbon;

use \Datetime;
use \DateTimeZone;
use \DateInterval;

use Jyotish\Ganita\Astro;
use Jyotish\Ganita\Ayanamsha;
use Jyotish\Muhurta\Hora;

include __DIR__ . "/config.php";


class Lib
{
    
    public $grahas;
    public $lagnas;

    public function __construct(){}

    
    public function calculator($params = [])
    {
        $latitude = $params['latitude'] ?? null;
        $longitude = $params['longitude'] ?? null;
        $year = $params['year'] ?? null;
        $month = $params['month'] ?? null;
        $day = $params['day'] ?? null;
        $hour = $params['hour'] ?? null;
        $min = $params['min'] ?? null;
        $sec = $params['sec'] ?? null;
        $time_zone = $params['time_zone'] ?? 'Asia/Tehran';
        $dst_hour = $params['dst_hour'] ?? 0;
        $dst_min = $params['dst_min'] ?? 0;
        $nesting = $params['nesting'] ?? 0;
        $varga = $params['varga'] ?? ["D1"];
        $infolevel = $params['infolevel'] ?? [];
        
        return $this->calculateChart(
            $latitude,
            $longitude,
            $year,
            $month,
            $day,
            $hour,
            $min,
            $sec,
            $time_zone,
            $dst_hour,
            $dst_min,
            $varga,
            $nesting,
            $infolevel
        );
    }

    public function calculateNow($latitude = '35.708309', $longitude = '51.380730')
    {
        $tz = $this->getNearestTimezone($latitude, $longitude);
        $now = new DateTime('now', new DateTimeZone($tz[0]));
        $year = $now->format('Y');
        $month = $now->format('m');
        $day = $now->format('d');
        $hour = $now->format('H');
        $min = $now->format('i');
        $sec = $now->format('s');
        $time_zone = $tz[1];
        $dst_hour = 0;
        $dst_min = 0;
        $nesting = 2;
        $varga = ["D1", "D9"];
        $infolevel = [];

        return $this->calculateChart(
            $latitude,
            $longitude,
            $year,
            $month,
            $day,
            $hour,
            $min,
            $sec,
            $time_zone,
            $dst_hour,
            $dst_min,
            $varga,
            $nesting,
            $infolevel
        );
    }

    public function getNearestTimezone($cur_lat, $cur_long, $country_code = null)
    {
        if ($country_code == 'ir') {

            $formatted = '+03:30';
            return ['Asia/Tehran', $formatted];
        }
        // $country_code = strtoupper($country_code);

        $timezone_ids = ($country_code) ? DateTimeZone::listIdentifiers(DateTimeZone::PER_COUNTRY, strtoupper($country_code))
            : DateTimeZone::listIdentifiers();

        // print_r($timezone_ids);
        if ($timezone_ids && is_array($timezone_ids) && isset($timezone_ids[0])) {

            $time_zone = '';
            $tz_distance = 0;

            //only one identifier?
            if (count($timezone_ids) == 1) {
                $time_zone = $timezone_ids[0];
            } else {

                foreach ($timezone_ids as $timezone_id) {
                    $timezone = new DateTimeZone($timezone_id);
                    $location = $timezone->getLocation();
                    $tz_lat   = $location['latitude'];
                    $tz_long  = $location['longitude'];
                    $tz_offset = $timezone->getOffset(new DateTime());

                    $theta    = $cur_long - $tz_long;
                    $distance = (sin(deg2rad($cur_lat)) * sin(deg2rad($tz_lat)))
                        + (cos(deg2rad($cur_lat)) * cos(deg2rad($tz_lat)) * cos(deg2rad($theta)));
                    $distance = acos($distance);
                    $distance = abs(rad2deg($distance));
                    // echo '<br />'.$timezone_id.' '.$distance; 

                    if (!$time_zone || $tz_distance > $distance) {
                        $time_zone   = $timezone_id;
                        $tz_distance = $distance;
                    }
                }
            }
            $tz = new DateTimeZone($time_zone);
            $date = new DateTime('now', $tz);

            // create a new date offset by the timezone offset
            // gets the interval as hours & minutes
            $offset = $tz->getOffset($date) . ' seconds';
            $dateOffset = clone $date;
            $dateOffset->sub(DateInterval::createFromDateString($offset));

            $interval = $dateOffset->diff($date);
            $formatted = $interval->format('%R%H:%I');


            return [$time_zone, $formatted];
        }
    }

    /*
     * $infolevel = ["basic", "ashtakavarga", "grahabala", "rashibala", "yogas", "panchanga", "transit"]
     *
     */
    public function calculateChart(
        $latitude,
        $longitude,
        $year,
        $month,
        $day,
        $hour,
        $min,
        $sec,
        $time_zone = 'Asia/Tehran',
        $dst_hour = 0,
        $dst_min = 0,
        $vargas = [
            "D1",
            "D2",
            "D3",
            "D4",
            "D7",
            "D9",
            "D10",
            "D12",
            "D16",
            "D20",
            "D24",
            "D27",
            "D30",
            "D40",
            "D45",
            "D60",
        ],
        $nesting = 4,
        array $infolevel = ["basic", "ashtakavarga", "grahabala", "rashibala", "yogas", "panchanga", "transit"]
    ) {

        $locality = new Locality([
            'longitude' => $longitude,
            'latitude' => $latitude,
            'altitude' => 0
        ]);

        # Use the caller's "+HH:MM" offset for the birth moment. Nearest-timezone guessing picked the
        # wrong zone for much of India (Mumbai -> Karachi, Delhi -> Kathmandu) and used today's DST.
        # format datetime for DateTime Object
        $datetime = sprintf("%s-%s-%s %s:%s:%s%s", $year, $month, $day, $hour, $min, $sec, $time_zone);
        $date = new DateTime($datetime);

        # perform DST offest
        $date->modify(sprintf("-%s hours", $dst_hour));
        $date->modify(sprintf("-%s minutes", $dst_min));

        # setup ephemeris and calculations
        $ganita = new Swetest(["swetest" => SWETEST_PATH]);
        $data = new Data($date, $locality, $ganita);
        $data->calcVargaData($vargas);
        $data->calcParams();
        if (in_array('panchanga', $infolevel)) {
            $data->calcPanchanga();
            $data->calcRising();
            $data->calcHora();
            $data->calcHora(Hora::TYPE_YAMA);
        }
        
        $analysis = new Analysis($data);
        $vargaData = $analysis->getVargaData('D1');


        if (in_array('ashtakavarga', $infolevel)) {
            // AshtakaVarga calculation
            $ashtakaVarga = new AshtakaVarga($data);
            // $vargaData['ashtakavarga'] = $ashtakaVarga->getSarvAshtakavarga(true);
            $vargaData['ashtakavarga'] = $ashtakaVarga->getBhinnAshtakavarga();
        }

        if (in_array('ayanamsa', $infolevel)) {
            $vargaData['ayanamsa'] = Ayanamsha::getAyanamsha();
        }

        if (in_array('grahabala', $infolevel)) {
            // GrahaBala calculation
            $grahaBala = new GrahaBala($data);
            $vargaData['grahabala'] = $grahaBala->getBala();

        }
        if (in_array('rashibala', $infolevel)) {
            // RashiBala calculation
            $rashiBala = new RashiBala($data);
            $vargaData['rashibala'] = $rashiBala->getBala();
        }


        $angaDefiner = new AngaDefiner($data);

        foreach ($vargaData['graha'] as $grahaKey => $value) {
            $nakshatra = $angaDefiner->getNakshatra(false, false, $grahaKey);
            $vargaData['graha'][$grahaKey]['nakshatra'] = $nakshatra;


            if (in_array('basic', $infolevel)) {
                $Graha = Graha::getInstance($grahaKey)->setEnvironment($data);
                $vargaData['graha'][$grahaKey]['astangata'] = $Graha->isAstangata(); // combustion
                $vargaData['graha'][$grahaKey]['rashiAvastha'] = $Graha->getRashiAvastha(); // dignity
                $vargaData['graha'][$grahaKey]['vargottama'] = $Graha->isVargottama(); // Vargottama
                $vargaData['graha'][$grahaKey]['yuddha'] = $Graha->isYuddha(); // graha is in planetary war
            }


            if (in_array('panchanga', $infolevel)) {
                $Graha = Graha::getInstance($grahaKey)->setEnvironment($data);
                $vargaData['graha'][$grahaKey]['astangata'] = $Graha->isAstangata(); // combustion
                $vargaData['graha'][$grahaKey]['rashiAvastha'] = $Graha->getRashiAvastha(); // dignity
                $vargaData['graha'][$grahaKey]['vargottama'] = $Graha->isVargottama(); // Vargottama
                $vargaData['graha'][$grahaKey]['yuddha'] = $Graha->isYuddha(); // graha is in planetary war
                $vargaData['graha'][$grahaKey]['gocharastha'] = $Graha->isGocharastha(); // gocharastha
                $vargaData['graha'][$grahaKey]['bhavaCharacter'] = $Graha->getBhavaCharacter(); // Bhava Character
                $vargaData['graha'][$grahaKey]['tempRelation'] = $Graha->getTempRelation(); // Get tatkalika (temporary) relations
                $vargaData['graha'][$grahaKey]['relation'] = $Graha->getRelation(); // Get summary relations
                $vargaData['graha'][$grahaKey]['yogakaraka'] = $Graha->isYogakaraka(); // yogakaraka
                $vargaData['graha'][$grahaKey]['mrityu'] = $Graha->isMrityu(); // graha is in mrityu bhaga
                $vargaData['graha'][$grahaKey]['pushkaraNavamsha'] = $Graha->isPushkara(Graha::PUSHKARA_NAVAMSHA); // graha is in pushkara navamsha
                $vargaData['graha'][$grahaKey]['pushkaraBhaga'] = $Graha->isPushkara(Graha::PUSHKARA_BHAGA); // graha is in pushkara bhaga
                $vargaData['graha'][$grahaKey]['avastha'] = $Graha->getAvastha(); // Get avastha of graha
                $vargaData['graha'][$grahaKey]['dispositor'] = $Graha->getDispositor(); // Get ruler of the bhava, where graha is positioned
            }
        }

        $nakshatra = $angaDefiner->getNakshatra(false, false, Lagna::KEY_LG);
        $vargaData['lagna'][Lagna::KEY_LG]['nakshatra'] = $nakshatra;

        $data->calcDasha(Dasha::TYPE_VIMSHOTTARI, null, ['nesting' => $nesting]);
        $dasha = $data->getData();


        if (in_array('panchanga', $infolevel)) {
            $vargaData['panchanga'] = $dasha['panchanga'];
        }

        $vargaData['dasha'] = $dasha['dasha']['vimshottari'];

        if (in_array('yogas', $infolevel)) {
            // Yoga calculations
            $data->calcYoga([
                Yoga::TYPE_DHANA,
                Yoga::TYPE_MAHAPURUSHA,
                Yoga::TYPE_NABHASHA,
                Yoga::TYPE_PARIVARTHANA,
                Yoga::TYPE_RAJA,
                Yoga::TYPE_SANNYASA,
                // Yoga::INTERPLAY_PARIVARTHANA,
                // Yoga::INTERPLAY_CONJUNCT,
                // Yoga::INTERPLAY_ASPECT
            ]);
            $yogas = $data->getData(['yoga']);
            $vargaData['yogas'] = $yogas;
        }

        return $vargaData;
    }


}
