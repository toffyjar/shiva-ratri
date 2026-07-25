import React from 'react';
import { View, StyleSheet, ActivityIndicator, Dimensions } from 'react-native';
import { Appbar } from 'react-native-paper';
import { Picker } from '@react-native-picker/picker';
import DrawChart from '../components/drawChart';
import HttpService from '../services/httpService';
import { formatChartData, formatTimeZoneForApi, ChartResponse, ChartProps, fallbackChartData, isChartData, isKundliAlreadySaved, SavedKundaliRecord } from '../utils/chartPage';
import { useKundliStore } from '../store/kundliStore';
import { useAuthStore } from '../store/authStore';
import { useSavedKundalisStore, SavedKundali } from '../store/savedKundalisStore';
import { SvgXml } from 'react-native-svg';

const { width } = Dimensions.get('window');

const apiService = new HttpService('http://192.168.1.10:9393');
const localApiService = new HttpService('http://192.168.1.10:5656');

const VARGA_OPTIONS = ['D1', 'D2', 'D3', 'D4', 'D7', 'D9', 'D10', 'D12', 'D16', 'D20', 'D24', 'D27', 'D30', 'D40', 'D45', 'D60'];

export default function Chart({ onBack }: ChartProps) {
  const { name, day, month, year, hour, minute, birthPlace, latitude, longitude, timeZone, isFemale } = useKundliStore();
  const token = useAuthStore((state) => state.token);
  const savedKundalis = useSavedKundalisStore((state) => state.kundalis);
  const addKundali = useSavedKundalisStore((state) => state.addKundali);
  const [chartData, setChartData] = React.useState<
    Record<number, { rashi: string; planets: string[] }>
  >(fallbackChartData);
  const [loading, setLoading] = React.useState(true);
  const [selectedVarga, setSelectedVarga] = React.useState('D1');
  const [svgString, setSvgString] = React.useState(`
    <svg height="100" width="100">
      <circle cx="50" cy="50" r="40" fill="#FF0000" />
    </svg>
  `);

  React.useEffect(() => {
    let cancelled = false;
    setLoading(true);

    (async () => {
      try {
        try {
          const currentKundli: SavedKundaliRecord = {
            name, day, month, year, hour, minute, birthPlace, latitude, longitude, timeZone, isFemale,
          };

          if (!isKundliAlreadySaved(savedKundalis, currentKundli)) {
            const saveResponse = await localApiService.post<SavedKundali>(
              '/api/v1/save-kundali',
              currentKundli,
              token ? { Authorization: `Bearer ${token}` } : undefined
            );

            if (saveResponse.ok && saveResponse.data) {
              addKundali(saveResponse.data);
            }
          }
        } catch (error) {
          console.log('Error saving kundali:', error);
        }

        const params = {
          latitude,
          longitude,
          year,
          month: month.padStart(2, '0'),
          day: day.padStart(2, '0'),
          hour,
          min: minute,
          sec: '0',
          time_zone: formatTimeZoneForApi(timeZone),
          dst_hour: '0',
          dst_min: '0',
          nesting: '0',
          size: Math.round((75 / 100) * width).toString(),
          varga: selectedVarga,
          // infolevel: 'basic,panchanga,transit',
        };
        const response = await apiService.getSvg('/api/chart/svg', params);
        // const formattedData = formatChartData(response.data.chart);
        if (response.ok && response?.data?.length > 0) {
          setSvgString(response.data);
          // setChartData(formattedData);
          // console.log('Chart data fetched successfully:', formattedData);
        }

        if (!cancelled) setLoading(false);

        // console.log('Fetching chart data:', response.data);
      } catch (error) {
        console.log('Error fetching chart data:', error);
        setLoading(false);
      }
    })();

    return () => {
      cancelled = true;
    };
  }, [name, day, month, year, hour, minute, birthPlace, latitude, longitude, timeZone, isFemale, selectedVarga, token, savedKundalis]);

  if (loading) {
    return (
      <View style={styles.wrapper}>
        <Appbar.Header style={{ backgroundColor: '#2196f3' }}>
          <Appbar.BackAction onPress={onBack} iconColor="#000" />
          <Appbar.Content title="Birth Chart" titleStyle={{ color: '#000' }} />
        </Appbar.Header>
        <View style={styles.pickerContainer}>
          <Picker
            selectedValue={selectedVarga}
            onValueChange={setSelectedVarga}
            mode="dropdown"
          >
            {VARGA_OPTIONS.map((v) => (
              <Picker.Item key={v} label={v} value={v} />
            ))}
          </Picker>
        </View>
        <View style={styles.container}>
          <ActivityIndicator size="large" />
        </View>
      </View>
    );
  }

  return (
    <View style={styles.wrapper}>
      <Appbar.Header  style={{ backgroundColor: '#2196f3'}}>
        <Appbar.BackAction onPress={onBack} iconColor="#000" />
        <Appbar.Content title="Back" titleStyle={{ color: '#000' }} />
      </Appbar.Header>
      <View style={styles.pickerContainer}>
        <Picker
          selectedValue={selectedVarga}
          onValueChange={setSelectedVarga}
          mode="dropdown"
        >
          {VARGA_OPTIONS.map((v) => (
            <Picker.Item key={v} label={v} value={v} />
          ))}
        </Picker>
      </View>
      <View style={styles.container}>
        <SvgXml xml={svgString} style={styles.container} />
        {/* <DrawChart chartData={chartData} /> */}
      </View>
    </View>
  );
}

const styles = StyleSheet.create({
  wrapper: { flex: 1 },
  container: { flex: 1, justifyContent: 'center', alignItems: 'center', backgroundColor: '#fff' },
  pickerContainer: { backgroundColor: '#fff', borderBottomWidth: 1, borderBottomColor: '#e0e0e0' },
});
